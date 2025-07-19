package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	crand "crypto/rand"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/golang-jwt/jwt"
	bolt "go.etcd.io/bbolt"
)

func FindOrCreateSchool(db *bolt.DB, clock Clock, name string, city string, zip int) (id string, err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		schools, err := tx.CreateBucketIfNotExists([]byte("schools"))
		if err != nil {
			return fmt.Errorf("cannot create 'schools': %v", err)
		}

		c := schools.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			school := schools.Bucket(k)

			if string(school.Get([]byte("name"))) != name {
				continue
			}
			if string(school.Get([]byte("city"))) != city {
				continue
			}
			if btoi(school.Get([]byte("zip"))) != zip {
				continue
			}
			id = string(k)
			return nil
		}

		id = RandomString(10)
		school, err := schools.CreateBucket([]byte(id))
		if err != nil {
			return fmt.Errorf("cannot create bucket for new school: %v", err)
		}

		err = school.Put([]byte(KeyName), []byte(name))
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeyCity), []byte(city))
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeyZip), itob(zip))
		if err != nil {
			return err
		}

		addCodes, _ := constSlice()
		addCode := randomWords(2, 100, addCodes)
		free, err := freeAddCodeRx(tx, addCode)
		if err != nil {
			return err
		}
		for !free {
			addCode = randomWords(2, 100, addCodes)
			free, err = freeAddCodeRx(tx, addCode)
			if err != nil {
				return err
			}
		}

		err = school.Put([]byte(KeyAddCode), []byte(addCode))
		if err != nil {
			return err
		}

		endTime := clock.Now().Add(time.Hour * 72).Truncate(time.Second).Format(time.RFC3339)

		err = school.Put([]byte(KeyRegEnd), []byte(endTime))
		if err != nil {
			return err
		}

		settings := openapi.Settings{
			Student2student: false,
			Lottery:         false,
			Odds:            500,
		}

		marshal, err := json.Marshal(settings)
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeySettings), marshal)
		if err != nil {
			return err
		}

		classes, err := school.CreateBucket([]byte(KeyClasses))
		if err != nil {
			return err
		}

		_, err = school.CreateBucket([]byte(KeyAuctions))
		if err != nil {
			return err
		}

		def := []struct {
			name   string
			period int
		}{
			{name: KeyFreshman, period: 9},
			{name: KeySophomores, period: 10},
			{name: KeyJuniors, period: 11},
			{name: KeySeniors, period: 12},
		}

		for _, c := range def {
			_, err = addClassDetailsTx(tx, classes, clock, c.name, c.period, true)
			if err != nil {
				return err
			}
		}

		cb, err := school.CreateBucket([]byte(KeyCB))
		if err != nil {
			return err
		}
		_, err = cb.CreateBucket([]byte(KeyAccounts))
		if err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func schoolsByZip(db *bolt.DB, zip int32) ([]openapi.ResponseSchoolsInner, error) {
	res := make([]openapi.ResponseSchoolsInner, 0)

	err := db.View(func(tx *bolt.Tx) error {
		schools := tx.Bucket([]byte("schools"))
		if schools == nil {
			return nil
		}

		c := schools.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			school := schools.Bucket(k)

			if btoi32(school.Get([]byte("zip"))) != zip {
				continue
			}
			schoolId := string(k)
			schoolName := string(school.Get([]byte("name")))
			res = append(res, openapi.ResponseSchoolsInner{
				Name: schoolName,
				Id:   schoolId,
			})
			return nil
		}
		return nil
	})

	return res, err
}

// opens a db.view to pass to getSchoolUsersRx
func getSchoolUsers(db *bolt.DB, schoolId string) (resp []openapi.UserNoHistory, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getSchoolUsersRx(tx, schoolId)
		return err
	})

	return
}

func makeToken(jwtSvc *token.Service, target UserInfo) (string, string, error) {

	tgt := token.User{
		Name: target.Name,
	}

	const ttl = 15 * time.Minute

	// XSRF = random 32-byte URL-safe string
	xsrfBytes := make([]byte, 32)
	if _, err := crand.Read(xsrfBytes); err != nil {
		return "", "", err
	}
	xsrf := base64.RawURLEncoding.EncodeToString(xsrfBytes)

	claims := token.Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        xsrf,                       // XSRF token goes into jti / id
			ExpiresAt: time.Now().Add(ttl).Unix(), // custom TTL, e.g. 15 min
			Issuer:    jwtSvc.Opts.Issuer,         // keep same issuer
		},
		User: &tgt,
	}

	jwtStr, err := jwtSvc.Token(claims) // sign only, no cookies written
	if err != nil {
		return "", "", err
	}
	return jwtStr, xsrf, nil
}

// gets all the users in a school
func getSchoolUsersRx(tx *bolt.Tx, schoolId string) (resp []openapi.UserNoHistory, err error) {
	school, err := schoolByIdTx(tx, schoolId)
	if err != nil {
		return
	}

	users := make([]string, 0)

	adminsBucket := school.Bucket([]byte(KeyAdmins))
	c := adminsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	teacherBucket := school.Bucket([]byte(KeyTeachers))
	c = teacherBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	studentsBucket := school.Bucket([]byte(KeyStudents))
	c = studentsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	//get all the users in the school
	for _, user := range users {
		userInfo, err := getUserInLocalStoreTx(tx, user)
		if err != nil {
			continue
		}

		rank := userInfo.Rank
		if rank == 0 {
			rank = 99999 // put zero-ranked users at the end
		}

		resp = append(resp, openapi.UserNoHistory{
			Id:            userInfo.Name,
			FirstName:     userInfo.FirstName,
			LastName:      userInfo.LastName,
			Role:          userInfo.Role,
			Rank:          rank, //if rank is 0 set it to 99999 so it is at the bottom of the list, this is to fix the sorting issue
			TaxableIncome: userInfo.TaxableIncome,
			Income:        userInfo.Income,
			College:       userInfo.College,
			NetWorth:      userInfo.NetWorth,
			LottoWin:      userInfo.LottoWin,
			LottoPlay:     userInfo.LottoPlay,
		})
	}

	return

}

// opens a db.view to pass to getSchoolsRx
func getSchools(db *bolt.DB) (resp []openapi.ResponseSchools, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getSchoolsRx(tx)
		return err
	})

	return
}

// opens a db.view to pass to getSchoolsRx
func getSchoolsRx(tx *bolt.Tx) (resp []openapi.ResponseSchools, err error) {
	schools := tx.Bucket([]byte("schools"))
	if schools == nil {
		return nil, fmt.Errorf("no schools available")
	}

	c := schools.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}

		school := schools.Bucket(k)

		item := openapi.ResponseSchools{
			Id:   string(k),
			Name: string(school.Get([]byte(KeyName))),
			City: string(school.Get([]byte(KeyCity))),
			Zip:  btoi32(school.Get([]byte(KeyZip))),
		}

		//loop through each bucket and count how many keys are in each
		adminsBucket := school.Bucket([]byte(KeyAdmins))
		c := adminsBucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Staff++
		}

		teacherBucket := school.Bucket([]byte(KeyTeachers))
		c = teacherBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Staff++
		}

		studentsBucket := school.Bucket([]byte(KeyStudents))
		c = studentsBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Students++
		}

		resp = append(resp, item)
	}
	return
}

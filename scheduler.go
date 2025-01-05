package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

// schedule a function to be called every minute
// returns channel to stop repeating
// to set up
//  runEveryMinute(func(time2 time.Time) {
//	  lgr.Printf("INFO Tick2 at", time2)
//   })
// to stop:   done <- true

func runEveryMinute(db *bolt.DB) (done chan bool) {
	ticker := time.NewTicker(60 * time.Second)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				coinGecko(db)
			}
		}
	}()

	return done
}

func runEveryDay(db *bolt.DB) (done chan bool) {
	ticker := time.NewTicker(24 * time.Hour)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				lgr.Printf("INFO running networths")
				schoolsNetworth(db)
			}
		}
	}()

	return done
}

func schoolsNetworthTx(tx *bolt.Tx) (err error) {
	schools := tx.Bucket([]byte(KeySchools))
	users := tx.Bucket([]byte(KeyUsers))

	c := schools.Cursor()

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		school := schools.Bucket(k)

		students := school.Bucket([]byte(KeyStudents))
		if students == nil {
			return fmt.Errorf("cannot find students bucket")
		}
		InnerCursor := students.Cursor()

		var resp []openapi.UserNoHistory

		for innerK, _ := InnerCursor.First(); innerK != nil; innerK, _ = InnerCursor.Next() {
			studentData := users.Get([]byte(innerK))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", innerK)
				continue
			}
			if student.Role != UserRoleStudent {
				lgr.Printf("ERROR student %s has role %d", innerK, student.Role)
				continue
			}

			nWorth, _ := StudentNetWorthTx(tx, student.Name).Float64()
			nUser := openapi.UserNoHistory{
				Id:        student.Email,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Rank:      student.Rank,
				NetWorth:  float32(nWorth),
			}

			resp = append(resp, nUser)

		}

		sort.SliceStable(resp, func(i, j int) bool {
			return resp[i].NetWorth > resp[j].NetWorth
		})

		for i := 0; i < len(resp); i++ {
			resp[i].Rank = int32(i + 1)
		}

		_, err = saveRanksTx(tx, resp)
		if err != nil {
			return fmt.Errorf("ERROR saving students ranks: %v", err)
		}

	}

	return nil

}

func schoolsNetworth(db *bolt.DB) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		return schoolsNetworthTx(tx)
	})

	if err != nil {
		lgr.Printf("ERROR networths are not updating")
	}

	return

}

func coinGecko(db *bolt.DB) (err error) {
	u, _ := url.ParseRequestURI("https://api.coingecko.com/api/v3/simple/price")
	q := u.Query()
	q.Set("ids", KeyCoins)
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	var decodedResp map[string]map[string]float32
	err = json.Unmarshal(body, &decodedResp)
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cryptosBucket, err := tx.CreateBucketIfNotExists([]byte(KeyCryptos))
		if err != nil {
			lgr.Printf(err.Error())
			return err
		}
		var cryptoInfo openapi.CryptoCb

		for k, v := range decodedResp {
			cryptoInfo.Usd = v["usd"]
			cryptoInfo.UpdatedAt = time.Now().Truncate(time.Second)
			marshal, err := json.Marshal(cryptoInfo)
			if err != nil {
				lgr.Printf(err.Error())
				return err
			}

			err = cryptosBucket.Put([]byte(k), marshal)
			if err != nil {
				lgr.Printf(err.Error())
				return err
			}
		}

		return nil
	})

	return

}

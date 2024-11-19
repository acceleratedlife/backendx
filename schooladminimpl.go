package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

func validateUserData(body UserInfo) (bool, error) {
	if body.SchoolId == "" {
		return false, fmt.Errorf("no school specified")
	}
	if body.Name == "" {
		return false, fmt.Errorf(" username is mandatory")
	}
	if len(body.Email) < 3 {
		return false, fmt.Errorf("wrong email format")
	}

	if body.FirstName == "" {
		return false, fmt.Errorf("first name is mandatory")
	}

	if body.LastName == "" {
		return false, fmt.Errorf("last name is mandatory")
	}

	return true, nil
}

func CreateSchoolAdmin(db *bolt.DB, body UserInfo) (openapi.ImplResponse, error) {

	_, err := validateUserData(body)
	if err != nil {
		return openapi.ImplResponse{}, err
	}

	if body.Role != UserRoleAdmin {
		return openapi.ImplResponse{}, fmt.Errorf("not a school admin")
	}
	v := openapi.ImplResponse{}

	err = db.Update(func(tx *bolt.Tx) error {

		newUser := UserInfo{
			Name:        body.Name,
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			Email:       body.Email,
			Confirmed:   true,
			PasswordSha: body.PasswordSha,
			SchoolId:    body.SchoolId,
			Role:        UserRoleAdmin,
		}

		err = AddUserTx(tx, newUser)
		if err != nil {
			return err
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(body.SchoolId))

		if school == nil {
			return fmt.Errorf("school is not created")
		}

		admins, err := school.CreateBucketIfNotExists([]byte(KeyAdmins))
		if err != nil {
			return err
		}

		_, err = school.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}

		admin := admins.Get([]byte(body.Name))

		if admin != nil {
			lgr.Printf("INFO admin %s is already registered ", body.Email)
			return nil
		}

		err = admins.Put([]byte(body.Name), []byte(""))
		if err != nil {
			return err
		}

		return nil
	})

	return v, err
}

func createTeacher(db *bolt.DB, newUser UserInfo) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {

		err = AddUserTx(tx, newUser)
		if err != nil {
			return err
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(newUser.SchoolId))

		if school == nil {
			return fmt.Errorf("school not found")
		}
		teachers, err := school.CreateBucketIfNotExists([]byte(KeyTeachers))
		if err != nil {
			return err
		}
		teacher, err := teachers.CreateBucket([]byte(newUser.Email))
		if err != nil {
			return err
		}
		_, err = teacher.CreateBucket([]byte(KeyClasses))
		if err != nil {
			return nil
		}
		return nil
	})
	return
}

func createStudent(db *bolt.DB, newUser UserInfo, pathId PathId) (err error) {
	jobDetails := getJob(db, KeyJobs, newUser.Job)
	max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
	min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
	diff := max.Sub(min)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := decimal.NewFromFloat32(r.Float32())
	newUser.Income = float32(random.Mul(diff).Add(min).Floor().InexactFloat64())

	newUser.CareerTransition = false
	err = db.Update(func(tx *bolt.Tx) error {
		err = AddUserTx(tx, newUser)
		if err != nil {
			// AddUser can fail if new user exists in the table as a student
			users := tx.Bucket([]byte(KeyUsers))
			if users == nil {
				return err
			}
			user := users.Get(userNameToB(newUser.Email))
			if user == nil {
				return err
			}
			var existingUser UserInfo

			err = json.Unmarshal(user, &existingUser)
			if err != nil {
				return err
			}
			if existingUser.Role != UserRoleStudent {
				return fmt.Errorf("you cannot add existing user with the role %d as a student to a class", existingUser.Role)
			}
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(newUser.SchoolId))

		if school == nil {
			return fmt.Errorf("school not found")
		}
		schoolStudentsBucket, err := school.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}
		schoolStudentBucket, err := schoolStudentsBucket.CreateBucketIfNotExists([]byte(newUser.Email))
		if err != nil {
			return err
		}
		history := []openapi.History{}
		row1 := openapi.History{
			Date:     time.Now(),
			NetWorth: 0.0,
		}
		history = append(history, row1)
		marshal, _ := json.Marshal(history)
		schoolStudentBucket.Put([]byte(KeyHistory), marshal)
		teachers, err := school.CreateBucketIfNotExists([]byte(KeyTeachers))
		if err != nil {
			return err
		}

		teacher := teachers.Bucket([]byte(pathId.teacherId))
		if teacher == nil {
			return fmt.Errorf("teacher not found")
		}
		classes := teacher.Bucket([]byte(KeyClasses))
		if classes == nil {
			return fmt.Errorf("classes not found")
		}
		class := classes.Bucket([]byte(pathId.classId))
		if class == nil {
			return fmt.Errorf("class not found")
		}
		students, err := class.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}
		student := students.Get([]byte(newUser.Email))
		if student != nil {
			return fmt.Errorf("student is already in the class")
		}
		return students.Put([]byte(newUser.Email), []byte(""))
	})
	return err
}

func taxSchool(db *bolt.DB, clock Clock, userDetails UserInfo, taxRate int32) error {
	return db.Update(func(tx *bolt.Tx) error {
		return taxSchoolTx(tx, clock, userDetails, taxRate)
	})
}

func taxSchoolTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, taxRate int32) error {
	school, err := SchoolByIdTx(tx, userDetails.SchoolId)
	if err != nil {
		return err
	}

	students := school.Bucket([]byte(KeyStudents))
	if students == nil {
		return fmt.Errorf("cannot find students bucket")
	}

	c := students.Cursor()

	users := tx.Bucket([]byte(KeyUsers))

	if taxRate > 0 {
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			studentData := users.Get([]byte(k))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}
			if student.Role != UserRoleStudent {
				lgr.Printf("ERROR student %s has role %d", k, student.Role)
				continue
			}

			adjustedTaxRate := decimal.NewFromInt32(taxRate).Div(decimal.NewFromInt(100))

			charge := adjustedTaxRate.Mul(decimal.NewFromInt32(student.TaxableIncome))
			chargeStudentUbuckTx(tx, clock, student, charge.Abs(), "Income Tax at: "+strconv.Itoa(int(taxRate))+"%", false)

			student.TaxableIncome = 0
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(k), marshal)
			if err != nil {
				return err
			}

		}

	} else {
		taxes, err := getTaxSliceRx(tx, userDetails)
		if err != nil {
			return err
		}

		var meanTax, standardDevTax decimal.Decimal
		var realTBrackets []decimal.Decimal

		for i := 1; i < 20; i++ {
			taxes = taxFilter(taxes, Key_lower_percentile*float64(i))

			meanTax = decimal.Avg(taxes[0], taxes[1:]...)

			standardDevTax, err = getStandardDevTax(taxes, meanTax)
			if err != nil {
				return err
			}

			realTBrackets = realTaxBrackets(meanTax, standardDevTax)
			if realTBrackets[0].GreaterThan(decimal.NewFromInt32(1500)) {
				break
			}
		}

		summativeTaxes := summativeTaxes(realTBrackets)

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			studentData := users.Get([]byte(k))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}
			if student.Role != UserRoleStudent {
				lgr.Printf("ERROR student %s has role %d", k, student.Role)
				continue
			}

			zScore := (decimal.NewFromInt32(student.TaxableIncome).Sub(meanTax)).Div(standardDevTax)
			adjustedTaxRate, bracketPosition := progressiveTaxRate(zScore)
			var charge decimal.Decimal

			if bracketPosition < 1 {
				charge = adjustedTaxRate.Mul(decimal.NewFromInt32(student.TaxableIncome))
			} else {
				//realDiff calculating how far above the previous bracket you are
				realDiff := decimal.NewFromInt32(student.TaxableIncome).Sub(realTBrackets[bracketPosition-1])
				//collect all the tax for the brackets they already passed through
				summativeTax := summativeTaxes[bracketPosition-1]
				//add what they already pass through plus the real diff at their current tax rate
				charge = summativeTax.Add(realDiff.Mul(adjustedTaxRate))
			}

			chargeStudentUbuckTx(tx, clock, student, charge.Abs(), "Income Tax at: "+strconv.Itoa(int(adjustedTaxRate.Mul(decimal.NewFromInt32(100)).IntPart()))+"%", false)

			student.TaxableIncome = 0
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(k), marshal)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func taxBrackets(db *bolt.DB, userDetails UserInfo) (resp []openapi.ResponseTaxBracket, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = taxBracketsRx(tx, userDetails)
		return err
	})

	return
}

func taxBracketsRx(tx *bolt.Tx, userDetails UserInfo) (resp []openapi.ResponseTaxBracket, err error) {

	taxes, err := getTaxSliceRx(tx, userDetails)
	if err != nil {
		return
	}

	var meanTax, standardDevTax decimal.Decimal
	var brackets []decimal.Decimal

	for i := 1; i < 20; i++ {
		taxes = taxFilter(taxes, Key_lower_percentile*float64(i))

		meanTax = decimal.Avg(taxes[0], taxes[1:]...)

		standardDevTax, err = getStandardDevTax(taxes, meanTax)
		if err != nil {
			return resp, err
		}

		brackets = realTaxBrackets(meanTax, standardDevTax)
		if brackets[0].GreaterThan(decimal.NewFromInt32(1500)) {
			break
		}
	}

	rates := make([]float32, 0)
	rates = append(rates, KeyTax0)
	rates = append(rates, KeyTax1)
	rates = append(rates, KeyTax2)
	rates = append(rates, KeyTax3)
	rates = append(rates, KeyTax4)
	rates = append(rates, KeyTax5)
	rates = append(rates, KeyTax6)

	for i, bracket := range brackets {
		resp = append(resp, openapi.ResponseTaxBracket{
			Bracket: float32(bracket.IntPart()),
			Rate:    rates[i],
		})
	}

	return
}

// gives the tax you pay if you are completly through a bracket
// for example if your in the second bracket or 12% percent then you must pay all of the 10% + some of the 12%
func summativeTaxes(brackets []decimal.Decimal) (summativeTaxes []decimal.Decimal) {
	summativeTaxes = append(summativeTaxes, brackets[0].Mul(decimal.NewFromFloat32(KeyTax0)))
	diff := brackets[1].Sub(brackets[0])
	summativeTaxes = append(summativeTaxes, summativeTaxes[0].Add(diff.Mul(decimal.NewFromFloat32(KeyTax1))))
	diff = brackets[2].Sub(brackets[1])
	summativeTaxes = append(summativeTaxes, summativeTaxes[1].Add(diff.Mul(decimal.NewFromFloat32(KeyTax2))))
	diff = brackets[3].Sub(brackets[2])
	summativeTaxes = append(summativeTaxes, summativeTaxes[2].Add(diff.Mul(decimal.NewFromFloat32(KeyTax3))))
	diff = brackets[4].Sub(brackets[3])
	summativeTaxes = append(summativeTaxes, summativeTaxes[3].Add(diff.Mul(decimal.NewFromFloat32(KeyTax4))))
	diff = brackets[5].Sub(brackets[4])
	summativeTaxes = append(summativeTaxes, summativeTaxes[4].Add(diff.Mul(decimal.NewFromFloat32(KeyTax5))))
	return
}

// gives the actuall tax breakpoints
// for example up to 50 @ 10%, up to 70 @ 12%, up to 90 @ 22%,...
// It wont' give you the percents just [50,70,90,...]
func realTaxBrackets(mean, standardDev decimal.Decimal) (brackets []decimal.Decimal) {
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ0))))
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ1))))
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ2))))
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ3))))
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ4))))
	brackets = append(brackets, mean.Add(standardDev.Mul(decimal.NewFromFloat32(KeyTaxZ5))))
	return
}

// gives the tax rate you are at and the rank of the tax rate
// for example if you were the top player then you would be at 37% which is the 7th bracket or 6th index
func progressiveTaxRate(zScore decimal.Decimal) (decimal.Decimal, int) {
	if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ0)) {
		return decimal.NewFromFloat32(KeyTax0), 0
	} else if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ1)) {
		return decimal.NewFromFloat32(KeyTax1), 1
	} else if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ2)) {
		return decimal.NewFromFloat32(KeyTax2), 2
	} else if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ3)) {
		return decimal.NewFromFloat32(KeyTax3), 3
	} else if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ4)) {
		return decimal.NewFromFloat32(KeyTax4), 4
	} else if zScore.LessThan(decimal.NewFromFloat32(KeyTaxZ5)) {
		return decimal.NewFromFloat32(KeyTax5), 5
	} else {
		return decimal.NewFromFloat32(KeyTax6), 6
	}
}

func getStandardDevTax(taxes []decimal.Decimal, mean decimal.Decimal) (standardDev decimal.Decimal, err error) { //need to test
	if len(taxes) == 0 {
		return decimal.Zero, fmt.Errorf("error length is zero and: %v", err)
	}
	squaredSums := decimal.Zero

	for _, value := range taxes {
		squaredSums = squaredSums.Add((value.Sub(mean)).Pow(decimal.NewFromFloat(2)))
	}

	variance := squaredSums.Div(decimal.NewFromInt(int64(len(taxes))))
	if variance.IsZero() {
		return decimal.Zero, fmt.Errorf("error Variance is zero and: %v", err)
	}

	standardDev = decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
	if standardDev.IsZero() {
		return decimal.Zero, fmt.Errorf("error Standard Deviation is zero and: %v", err)
	}

	return
}

func getTaxSlice(db *bolt.DB, userDetails UserInfo) (taxes []decimal.Decimal, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		taxes, err = getTaxSliceRx(tx, userDetails)
		return err
	})

	return
}

func getTaxSliceRx(tx *bolt.Tx, userDetails UserInfo) (taxes []decimal.Decimal, err error) {
	school, err := SchoolByIdTx(tx, userDetails.SchoolId)
	if err != nil {
		return
	}

	students := school.Bucket([]byte(KeyStudents))
	if students == nil {
		return taxes, fmt.Errorf("cannot find students bucket")
	}

	c := students.Cursor()

	users := tx.Bucket([]byte(KeyUsers))

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		studentData := users.Get([]byte(k))
		var student UserInfo
		err = json.Unmarshal(studentData, &student)
		if err != nil {
			lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
			continue
		}
		if student.Role != UserRoleStudent {
			lgr.Printf("ERROR student %s has role %d", k, student.Role)
			continue
		}

		taxes = append(taxes, decimal.NewFromInt32(student.TaxableIncome))
	}

	sort.Slice(taxes, func(i, j int) bool {
		return taxes[i].LessThan(taxes[j])
	})

	return
}

func taxFilter(actualTaxes []decimal.Decimal, percentile float64) (adjustedTaxes []decimal.Decimal) {
	breakpoint := int(percentile * float64(len(actualTaxes)))
	for i, value := range actualTaxes {
		if i <= breakpoint {
			adjustedTaxes = append(adjustedTaxes, actualTaxes[breakpoint])
		} else {
			adjustedTaxes = append(adjustedTaxes, value)
		}
	}

	return
}

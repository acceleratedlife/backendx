package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestCreateStudent(t *testing.T) {

	db, dbTearDown := OpenTestDB("createStudent")
	defer dbTearDown()

	_, schools, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	user := UserInfo{
		Name:             students[0],
		CareerTransition: false,
		College:          false,
		FirstName:        "ss",
		LastName:         "ss",
		Email:            students[0],
		SchoolId:         schools[0],
		Role:             UserRoleTeacher,
	}

	path := PathId{
		schoolId:  schools[0],
		teacherId: teachers[0],
		classId:   classes[0],
	}

	err = createStudent(db, user, path)
	require.NotNil(t, err)
}

func TestCreateStudentIncome(t *testing.T) {

	db, dbTearDown := OpenTestDB("createStudentIncome")
	defer dbTearDown()

	_, schools, teachers, classes, _, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	user := UserInfo{
		Name:             "ss@ss.com",
		CareerTransition: false,
		College:          false,
		FirstName:        "ss",
		LastName:         "ss",
		Email:            "ss@ss.com",
		SchoolId:         schools[0],
		Role:             UserRoleStudent,
		Job:              getJobId(db, KeyJobs),
	}

	path := PathId{
		schoolId:  schools[0],
		teacherId: teachers[0],
		classId:   classes[0],
	}

	err = createStudent(db, user, path)
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, "ss@ss.com")
	require.Nil(t, err)
	require.GreaterOrEqual(t, student.Income, float32(104))
}

func TestTaxSchoolProgressive(t *testing.T) {

	db, dbTearDown := OpenTestDB("TaxSchoolProgressive")
	defer dbTearDown()
	clock := TestClock{}

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 3)
	require.Nil(t, err)

	var student UserInfo

	for _, d := range students {
		student, err = getUserInLocalStore(db, d)
		require.Nil(t, err)
		require.Equal(t, int32(0), student.TaxableIncome)
		clock.TickOne(time.Hour * 24)
		DailyPayIfNeeded(db, &clock, student)
		clock.TickOne(time.Hour * 24)
		DailyPayIfNeeded(db, &clock, student)
		clock.TickOne(time.Hour * 24)
		DailyPayIfNeeded(db, &clock, student)
		student, err = getUserInLocalStore(db, d)
		require.Nil(t, err)
		require.Greater(t, student.TaxableIncome, int32(0))
	}

	getSchoolStudents(db, student)
	err = taxSchool(db, &clock, student, 0)
	require.Nil(t, err)

	student, err = getUserInLocalStore(db, student.Email)
	require.Nil(t, err)
	require.Equal(t, int32(0), student.TaxableIncome)
	require.Greater(t, float64(student.Income*3), StudentNetWorth(db, student.Email).InexactFloat64())

}

func TestTaxSchoolFlat(t *testing.T) {

	db, dbTearDown := OpenTestDB("TaxSchoolFlat")
	defer dbTearDown()
	clock := TestClock{}

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 3)
	require.Nil(t, err)

	var student UserInfo

	for _, d := range students {
		student, err = getUserInLocalStore(db, d)
		require.Nil(t, err)
		require.Equal(t, int32(0), student.TaxableIncome)
		clock.TickOne(time.Hour * 24)
		DailyPayIfNeeded(db, &clock, student)
		student, err = getUserInLocalStore(db, d)
		require.Nil(t, err)
		require.Greater(t, student.TaxableIncome, int32(0))
	}

	getSchoolStudents(db, student)
	err = taxSchool(db, &clock, student, 11)
	require.Nil(t, err)

	student, err = getUserInLocalStore(db, student.Email)
	require.Nil(t, err)
	require.Equal(t, int32(0), student.TaxableIncome)
	require.Greater(t, float64(student.Income), StudentNetWorth(db, student.Email).InexactFloat64())

}

func TestGetMeanAndSDTax(t *testing.T) {

	db, dbTearDown := OpenTestDB("TaxSchoolFlat")
	defer dbTearDown()

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 10)
	require.Nil(t, err)

	income := []int32{251, 451, 851, 1651, 3251, 6451, 12851, 25651, 51251, 104851}

	err = db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		for i, d := range students {
			studentData := users.Get([]byte(d))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				return err
			}

			student.TaxableIncome = income[i]
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(d), marshal)
			if err != nil {
				return err
			}

		}
		return err
	})

	user, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	mean, err := getMeanTax(db, user)
	require.Nil(t, err)

	require.Equal(t, int64(20751), mean.IntPart())

	sd, err := getStandardDevTax(db, user, mean)
	require.Nil(t, err)
	require.Equal(t, int64(31927), sd.IntPart())

}

func TestTaxBrackets(t *testing.T) {

	db, dbTearDown := OpenTestDB("TaxBrackets")
	defer dbTearDown()

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 10)
	require.Nil(t, err)

	income := []int32{251, 451, 851, 1651, 3251, 6451, 12851, 25651, 51251, 104851}

	err = db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		for i, d := range students {
			studentData := users.Get([]byte(d))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				return err
			}

			student.TaxableIncome = income[i]
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(d), marshal)
			if err != nil {
				return err
			}

		}
		return err
	})

	user, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	resp, err := taxBrackets(db, user)
	require.Nil(t, err)

	require.Greater(t, resp[4].Bracket, float32(0))
	require.Greater(t, resp[1].Bracket, resp[0].Bracket)

}

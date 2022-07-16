package main

import (
	"testing"

	"github.com/stretchr/testify/require"
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

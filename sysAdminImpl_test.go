package main

import (
	"testing"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/require"
)

// If you just can spec then ResponseSchool model should not have omitEmpty on students
func TestSchoolByZip(t *testing.T) {

	lgr.Printf("INFO TestSchoolByZip")
	t.Log("INFO TestSchoolByZip")
	db, dbTearDown := OpenTestDB("refundCDFunction")
	defer dbTearDown()
	CreateTestAccounts(db, 3, 1, 1, 1)

	schools, err := schoolsByZip(db, 0) //there is 1 school with zip 0 but 0 will fail the if statement. No school in real server have zip 0
	require.Nil(t, err)

	require.Equal(t, 3, len(schools))

	schools, err = schoolsByZip(db, 1) //there exists 1 school with zip 1
	require.Nil(t, err)

	require.Equal(t, 1, len(schools))

}

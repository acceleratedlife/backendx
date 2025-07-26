package main

import (
	"testing"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

/*
	-----------------------------------------------------------------------
	  tiny helper – Bolt stores zip as []byte; we need an encoder in tests

------------------------------------------------------------------------
*/

/*
	-----------------------------------------------------------------------
	  schoolsByZip

------------------------------------------------------------------------
*/
func TestSchoolsByZip(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	_, schools, _, _, _, err := CreateTestAccounts(db, 3, 1, 1, 1)
	require.NoError(t, err)

	wantZIP, otherZIP := int32(90210), int32(10001)
	require.NoError(t, db.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte("schools"))
		// first school gets the target ZIP, others something else
		_ = sb.Bucket([]byte(schools[0])).Put([]byte("zip"), itob32(wantZIP))
		for _, id := range schools[1:] {
			_ = sb.Bucket([]byte(id)).Put([]byte("zip"), itob32(otherZIP))
		}
		return nil
	}))

	list, err := schoolsByZip(db, wantZIP)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, schools[0], list[0].Id)
}

/*
	-----------------------------------------------------------------------
	  getSchoolUsers / getSchoolUsersRx

------------------------------------------------------------------------
*/
func TestGetSchoolUsers(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, schools, teachers, _, students, err :=
		CreateTestAccounts(db, 1, 2, 3, 4) // 1 school, 2 admins, 3 teachers, 4 students
	require.NoError(t, err)

	users, err := getSchoolUsers(db, schools[0])
	require.NoError(t, err)

	expectTotal := len(admins) + len(teachers) + len(students)
	assert.Len(t, users, expectTotal)

	// a known student must be present
	present := false
	for _, u := range users {
		if u.Id == students[0] {
			present = true
			break
		}
	}
	assert.True(t, present)
}

/*
	-----------------------------------------------------------------------
	  makeToken – ensure JWT ↔︎ XSRF coherence

------------------------------------------------------------------------
*/
func TestMakeToken(t *testing.T) {
	// minimal auth/token service
	sr := token.SecretFunc(func(_ string) (string, error) { return "test_secret", nil })
	authSvc := auth.NewService(auth.Opts{
		SecretReader:   sr,
		TokenDuration:  time.Hour,
		CookieDuration: time.Hour,
		DisableXSRF:    true,
		Issuer:         "AL‑test",
	})

	target := UserInfo{Name: "alice@example.com"}
	jwtStr, xsrf, err := makeToken(authSvc.TokenService(), target)
	require.NoError(t, err)

	assert.NotEmpty(t, jwtStr)
	assert.NotEmpty(t, xsrf)

	claims, err := authSvc.TokenService().Parse(jwtStr)
	require.NoError(t, err)
	assert.Equal(t, xsrf, claims.Id)
	assert.Equal(t, target.Name, claims.User.Name)

	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
	assert.Greater(t, ttl, 14*time.Minute) // ≈ 15 min TTL as coded
	assert.LessOrEqual(t, ttl, 15*time.Minute)
}

/*
	-----------------------------------------------------------------------
	  getSchools / getSchoolsRx – counts of staff & students

------------------------------------------------------------------------
*/
func TestGetSchoolsCounts(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, schools, teachers, _, students, err :=
		CreateTestAccounts(db, 1, 2, 3, 4) // 1 school, 2 admins, 3 teachers, 4 students
	require.NoError(t, err)

	list, err := getSchools(db)
	require.NoError(t, err)
	require.Len(t, list, 1)

	s := list[0]
	assert.Equal(t, schools[0], s.Id)
	assert.Equal(t, int32(len(students)), s.Students)
	assert.Equal(t, int32(len(admins)+len(teachers)), s.Staff)
}

/*
	-----------------------------------------------------------------------
	  getSchools – error when no "schools" bucket

------------------------------------------------------------------------
*/
func TestGetSchools_NoBucket(t *testing.T) {
	db, closeDB := OpenTestDB("") // brand‑new DB, no schools created
	defer closeDB()

	_, err := getSchools(db)
	require.Error(t, err)
	assert.Equal(t, "no schools available", err.Error())
}

func TestMessageEmptyImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	_, _, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.NoError(t, err)

	err = message(db, nil, nil, nil, false, false)
	require.Error(t, err)
	assert.Equal(t, "no message sent", err.Error())
}

func TestMessageUserImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, &students[0], nil, false, false)
	require.NoError(t, err)

	user, err := getUserInLocalStore(db, students[0])
	require.NoError(t, err)
	assert.Len(t, user.Messages, 1)
	assert.Equal(t, theMessage, user.Messages[0])
}

func TestMessageAllImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, nil, true, true)
	require.NoError(t, err)

	users := append(admins, teachers...)
	users = append(users, students...)

	for _, id := range users {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 1)
		assert.Equal(t, theMessage, user.Messages[0])
	}
}

func TestMessageAllStaffImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, nil, false, true)
	require.NoError(t, err)

	for _, id := range admins {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 1)
		assert.Equal(t, theMessage, user.Messages[0])
	}

	for _, id := range teachers {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 1)
		assert.Equal(t, theMessage, user.Messages[0])
	}

	for _, id := range students {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}
}

func TestMessageAllStudentsImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, nil, true, false)
	require.NoError(t, err)

	for _, id := range admins {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}

	for _, id := range teachers {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}

	for _, id := range students {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 1)
		assert.Equal(t, theMessage, user.Messages[0])
	}
}

func TestMessageSchoolImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, schools, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, &schools[0], true, true)
	require.NoError(t, err)

	for _, id := range admins {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}

	for _, id := range teachers {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}

	for _, id := range students {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}
}

func TestMessageSchoolStaffImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, schools, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, &schools[0], false, true)
	require.NoError(t, err)

	for _, id := range admins {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}

	for _, id := range teachers {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}

	for _, id := range students {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}
}

func TestMessageSchoolStudentsImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	admins, schools, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 1, 5)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, nil, &schools[0], true, false)
	require.NoError(t, err)

	for _, id := range admins {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}

	for _, id := range teachers {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		assert.Len(t, user.Messages, 0)
	}

	for _, id := range students {
		user, err := getUserInLocalStore(db, id)
		require.NoError(t, err)
		if user.SchoolId == schools[0] {
			assert.Len(t, user.Messages, 1)
			assert.Equal(t, theMessage, user.Messages[0])
		} else {
			assert.Len(t, user.Messages, 0)
		}
	}
}

func TestDeleteSchool(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	_, schools, _, _, _, err := CreateTestAccounts(db, 3, 1, 1, 1)
	require.NoError(t, err)

	users, err := getSchoolUsers(db, schools[0])
	require.NoError(t, err)
	err = deleteSchool(db, schools[0])
	require.NoError(t, err)

	//verify that there are only 2 schools left
	remainingSchools, err := getSchools(db)
	require.NoError(t, err)
	assert.Len(t, remainingSchools, 2)
	//verify that the admin and their teacher account is gone

	_, err = getUserInLocalStore(db, users[0].Id)
	require.Error(t, err)
}

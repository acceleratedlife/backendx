package main

import (
	"fmt"
	"math/rand"

	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

// finds school in schools bucket by addCode
func schoolByAddCodeTx(schools *bolt.Bucket, addCode string) (*bolt.Bucket, string, error) {
	c := schools.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}

		school := schools.Bucket(k)
		if school == nil {
			lgr.Printf("ERROR school not found. bucket is nil")
			continue
		}

		addCodeTx := school.Get([]byte(KeyAddCode))
		if addCodeTx != nil &&
			string(addCodeTx) == addCode {
			return school, string(k), nil
		}
	}
	return nil, "", fmt.Errorf("school not found")
}

func SchoolByIdTx(tx *bolt.Tx, schoolId string) (school *bolt.Bucket, err error) {
	schools := tx.Bucket([]byte(KeySchools))
	if schools == nil {
		return nil, fmt.Errorf("no schools available")
	}

	school = schools.Bucket([]byte(schoolId))
	if school == nil {
		return nil, fmt.Errorf("school not found")
	}
	return school, nil

}

// RoleByAddCode determines role by addCode and school, for students it determines the class
func RoleByAddCode(db *bolt.DB, code string) (role int32, pathId PathId, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		schools := tx.Bucket([]byte(KeySchools))
		if schools == nil {
			return fmt.Errorf("schools not found")
		}

		c := schools.Cursor()

		for currentSchoolId, v := c.First(); currentSchoolId != nil; currentSchoolId, v = c.Next() {
			if v != nil {
				continue
			}

			school := schools.Bucket(currentSchoolId)
			if school == nil {
				lgr.Printf("ERROR school %s not found. bucket is nil", string(currentSchoolId))
				continue
			}

			addCodeTx := school.Get([]byte(KeyAddCode))
			if addCodeTx != nil &&
				string(addCodeTx) == code {
				role = UserRoleTeacher
				pathId.schoolId = string(currentSchoolId)
				return nil
			}

			schoolClasses := school.Bucket([]byte(KeyClasses))
			if schoolClasses != nil {
				c := schoolClasses.Cursor()
				for currentClassId, v := c.First(); currentClassId != nil; currentClassId, v = c.Next() {
					if v != nil {
						continue
					}
					class := schoolClasses.Bucket(currentClassId)
					if class == nil {
						return fmt.Errorf("class not found")
					}
					addCodeTx := class.Get([]byte(KeyAddCode))
					if addCodeTx != nil && string(addCodeTx) == code {
						admins := school.Bucket([]byte(KeyAdmins))
						cAdmins := admins.Cursor()
						adminId, _ := cAdmins.First()
						role = UserRoleStudent
						pathId.schoolId = string(currentSchoolId)
						pathId.classId = string(currentClassId)
						pathId.teacherId = string(adminId)
						return nil
					}
				}
			}

			teachers := school.Bucket([]byte(KeyTeachers))
			if teachers == nil {
				continue
			}

			res := teachers.ForEach(func(teacherId, v []byte) error {
				if v != nil {
					return nil
				}
				teacher := teachers.Bucket(teacherId)
				if teacher == nil {
					return nil
				}
				classesBucket := teacher.Bucket([]byte(KeyClasses))
				if classesBucket == nil {
					return nil
				}

				res := classesBucket.ForEach(func(currentClassId, v []byte) error {
					if v != nil {
						return nil
					}
					class := classesBucket.Bucket(currentClassId)
					if class == nil {
						return nil
					}
					addCodeTx := class.Get([]byte(KeyAddCode))
					if addCodeTx != nil && string(addCodeTx) == code {
						role = UserRoleStudent
						pathId.schoolId = string(currentSchoolId)
						pathId.teacherId = string(teacherId)
						pathId.classId = string(currentClassId)
						return fmt.Errorf("found")
					}
					return nil
				})

				return res
			})

			if res != nil && res.Error() == "found" {
				return nil
			}

		}
		return fmt.Errorf("Invalid Add Code")
	})

	return
}

func getJobId(db *bolt.DB, key string) (job string) {
	_ = db.View(func(tx *bolt.Tx) error {
		job = getJobIdRx(tx, key)
		return nil
	})

	return job
}

func getJobIdRx(tx *bolt.Tx, key string) (job string) {
	jobs := tx.Bucket([]byte(key))
	bucketStats := jobs.Stats()
	pick := rand.Intn(bucketStats.KeyN)
	c := jobs.Cursor()
	i := 0
	for k, _ := c.First(); k != nil && i <= pick; k, _ = c.Next() {
		if i != pick {
			i++
			continue
		}

		i++

		job = string(k)

	}

	return
}

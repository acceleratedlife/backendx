package main

import (
	"fmt"
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

		addCodeTx := school.Get([]byte("addCode"))
		if addCodeTx != nil &&
			string(addCodeTx) == addCode {
			return school, string(k), nil
		}
	}
	return nil, "", fmt.Errorf("school not found")
}

func SchoolByIdTx(tx *bolt.Tx, schoolId string) (school *bolt.Bucket, err error) {
	schools := tx.Bucket([]byte("schools"))
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
		schools := tx.Bucket([]byte("schools"))
		if schools == nil {
			return fmt.Errorf("school not found")
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

			addCodeTx := school.Get([]byte("addCode"))
			if addCodeTx != nil &&
				string(addCodeTx) == code {
				role = UserRoleTeacher
				pathId.schoolId = string(currentSchoolId)
				return nil
			}

			teachers := school.Bucket([]byte("teachers"))
			if teachers == nil {
				return nil
			}

			res := teachers.ForEach(func(teacherId, v []byte) error {
				if v != nil {
					return nil
				}
				teacher := teachers.Bucket(teacherId)
				if teacher == nil {
					return nil
				}

				res := teacher.ForEach(func(currentClassId, v []byte) error {
					if v != nil {
						return nil
					}
					class := teacher.Bucket(currentClassId)
					if class == nil {
						return nil
					}
					role = UserRoleStudent
					pathId.schoolId = string(currentSchoolId)
					pathId.teacherId = string(teacherId)
					pathId.classId = string(currentClassId)
					return fmt.Errorf("found")

				})

				return res
			})

			if res != nil && res.Error() == "found" {
				return nil
			}

		}
		return fmt.Errorf("school not found")
	})

	return
}

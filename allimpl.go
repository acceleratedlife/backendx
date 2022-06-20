package main

import (
	"encoding/json"
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

func getClassAtSchoolTx(tx *bolt.Tx, schoolId, classId string) (classBucket *bolt.Bucket, parentBucket *bolt.Bucket, err error) {

	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return nil, nil, err
	}

	classes := school.Bucket([]byte(KeyClasses))
	if classes != nil {
		classBucket := classes.Bucket([]byte(classId))
		if classBucket != nil {
			return classBucket, classes, nil
		}
	}

	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return nil, nil, fmt.Errorf("no teachers at school")
	}
	cTeachers := teachers.Cursor()
	for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
		if v != nil {
			continue
		}
		teacher := teachers.Bucket(k)
		if teacher == nil {
			continue
		}
		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			continue
		}
		classBucket = classesBucket.Bucket([]byte(classId)) //found the class
		if classBucket == nil {
			continue
		}
		return classBucket, classesBucket, nil
	}
	return nil, nil, fmt.Errorf("class not found")
}

func PopulateClassMembers(tx *bolt.Tx, classBucket *bolt.Bucket) (Members []openapi.ClassWithMembersMembers, err error) {
	Members = make([]openapi.ClassWithMembersMembers, 0)
	students := classBucket.Bucket([]byte(KeyStudents))
	if students == nil {
		return Members, nil
	}
	cStudents := students.Cursor()
	for k, _ := cStudents.First(); k != nil; k, _ = cStudents.Next() { //iterate students bucket
		user, err := getUserInLocalStoreTx(tx, string(k))
		if err != nil {
			return nil, err
		}

		nWorth, _ := StudentNetWorthTx(tx, user.Email).Float64()
		nUser := openapi.ClassWithMembersMembers{
			Id:        user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Rank:      user.Rank,
			NetWorth:  float32(nWorth),
		}
		Members = append(Members, nUser)
	}
	return Members, nil
}

func getStudentHistory(db *bolt.DB, userName string, schoolId string) (history []openapi.History, err error) {
	_ = db.View(func(tx *bolt.Tx) error {
		history, err = getStudentHistoryTX(tx, userName)
		return nil
	})
	return
}
func getStudentHistoryTX(tx *bolt.Tx, userName string) (history []openapi.History, err error) {
	student, err := getStudentBucketRoTx(tx, userName)
	if err == nil {
		return nil, err
	}

	historyData := student.Get([]byte(KeyHistory))
	if historyData == nil {
		return nil, fmt.Errorf("Failed to get history")
	}
	err = json.Unmarshal(historyData, &history)
	if err != nil {
		return nil, fmt.Errorf("ERROR cannot unmarshal History")
	}
	return
}

func getStudentAccountRx(tx *bolt.Tx, bAccount *bolt.Bucket, id string) (resp openapi.ResponseCurrencyExchange, err error) {
	// historyData := bAccount.Get([]byte(KeyHistory))
	// if historyData == nil {
	// 	return resp, fmt.Errorf("Failed to get history")
	// }
	// err = json.Unmarshal(historyData, &resp.History)
	// if err != nil {
	// 	return resp, fmt.Errorf("ERROR cannot unmarshal History")
	// }

	balanceData := bAccount.Get([]byte(KeyBalance))
	err = json.Unmarshal(balanceData, &resp.Balance)
	resp.Id = id

	return
}

func getCBaccountDetailsRx(tx *bolt.Tx, userDetails UserInfo, account openapi.ResponseCurrencyExchange) (finalAccount openapi.ResponseCurrencyExchange, err error) {
	cb, err := getCbRx(tx, userDetails.SchoolId)
	if err != nil {
		return
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return finalAccount, fmt.Errorf("cannot find cb buck accounts")
	}

	bAccount := accounts.Bucket([]byte(account.Id))
	if bAccount == nil {
		return finalAccount, fmt.Errorf("cannot find cb buck account")
	}

	// historyData := bAccount.Get([]byte(KeyHistory))
	// if historyData == nil {
	// 	return finalAccount, fmt.Errorf("Failed to get history")
	// }
	// err = json.Unmarshal(historyData, &account.History)
	// if err != nil {
	// 	return finalAccount, fmt.Errorf("ERROR cannot unmarshal History")
	// }

	finalAccount = account

	return
}

func getBuckNameTx(tx *bolt.Tx, id string) (string, error) {
	if id == KeyDebt {
		return "Debt", nil
	}

	if id == CurrencyUBuck {
		return "UBuck", nil
	}

	user, err := getUserInLocalStoreTx(tx, id)
	if err != nil {
		return "", err
	}

	return user.LastName + " Buck", nil
}

package main

import (
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

func UserByIdTx(tx *bolt.Tx, userId openapi.RequestUser) (user *bolt.Bucket, err error) {
	users := tx.Bucket([]byte("users"))
	if users == nil {
		return nil, fmt.Errorf("no users available")
	}

	user = users.Bucket([]byte(userId.Id))
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil

}

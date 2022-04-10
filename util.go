package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	bolt "go.etcd.io/bbolt"
)

func RandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func itob32(v int32) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi returns an 8-byte big endian representation of v.
func btoi(b []byte) int {

	u := binary.BigEndian.Uint64(b)
	return int(u)
}

// btoi32 returns an 8-byte big endian representation of v.
func btoi32(b []byte) int32 {

	u := binary.BigEndian.Uint64(b)
	return int32(u)
}

func userNameToB(username string) []byte {
	return []byte(username)
}

func EncodePassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("al"+password)))
}

// iterates over buckets inside given bucket
func iterateBuckets(mainBucket *bolt.Bucket, teacherHandler func(bucket *bolt.Bucket, key []byte)) {
	c := mainBucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}
		teacher := mainBucket.Bucket(k)
		teacherHandler(teacher, k)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"strconv"
)

func backUpHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		err := db.View(func(tx *bolt.Tx) error {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
			w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
			_, err := tx.WriteTo(w)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

type NewSchoolRequest struct {
	School    string
	FirstName string
	LastName  string
	Email     string
}

type NewSchoolResponse struct {
	AdminPassword string
}

func newSchoolHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := db.Update(func(tx *bolt.Tx) error {

			var request NewSchoolRequest
			defer r.Body.Close()
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&request)
			if err != nil {
				return fmt.Errorf("cannot parse request body: %v", err)
			}

			lgr.Printf("INFO new school request: %v", request)

			if request.Email == "" {
				return fmt.Errorf("email is mandatory")
			}
			if request.LastName == "" {
				return fmt.Errorf("lastname is mandatory")
			}

			/////////////////////////////
			// custom logic to create a new school

			response := NewSchoolResponse{
				AdminPassword: RandomString(8),
			}
			lgr.Printf("INFO new school NOT created")
			///////////////////////////

			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			err = encoder.Encode(response)
			if err != nil {
				lgr.Printf("ERROR failed to send")
			}

			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
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
	City      string
	Zip       int
}

type NewSchoolResponse struct {
	AdminPassword string
}

func newSchoolHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request NewSchoolRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new school request: %v", request)

		if request.Email == "" {
			err = fmt.Errorf("email is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if request.LastName == "" {
			err = fmt.Errorf("lastname is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		/////////////////////////////
		// custom logic to create a new school

		response := NewSchoolResponse{
			AdminPassword: RandomString(8),
		}
		err = createNewSchool(db, request, response.AdminPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new school %s created", request.School)
		///////////////////////////

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

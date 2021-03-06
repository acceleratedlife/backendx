package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	openapi "github.com/acceleratedlife/backend/go"
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

type resetPasswordRequest struct {
	Email string
}

type eventRequest struct {
	Positive    bool `json:",omitempty"`
	Description string
	Title       string `json:",omitempty"`
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

func resetPasswordHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request resetPasswordRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if request.Email == "" {
			err = fmt.Errorf("email is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := getUserInLocalStore(db, request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := resetPassword(db, user)

		lgr.Printf("Password reset for %s ", request.Email)

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func addJobsHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requests []Job
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requests)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new job request: %v", requests)

		for i, request := range requests {

			if request.Title == "" {
				err = fmt.Errorf("title is mandatory on: " + strconv.Itoa(i))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if request.Pay == 0 {
				err = fmt.Errorf("pay is mandatory on: " + strconv.Itoa(i))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if request.Description == "" {
				err = fmt.Errorf("description is mandatory on: " + strconv.Itoa(i))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			title := request.Title
			request.Title = ""

			if request.College {
				request.College = false
				marshal, err := json.Marshal(request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = createJobOrEvent(db, marshal, KeyCollegeJobs, title)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				marshal, err := json.Marshal(request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = createJobOrEvent(db, marshal, KeyJobs, title)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			lgr.Printf("job created for %s ", request.Title)

		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode("success")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func addEventsHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requests []eventRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requests)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new event request: %v", requests)

		for i, request := range requests {
			if request.Description == "" {
				err = fmt.Errorf("description is mandatory on: " + strconv.Itoa(i))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if request.Title == "" {
				err = fmt.Errorf("title is mandatory on: " + strconv.Itoa(i))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			title := request.Title
			request.Title = ""

			if request.Positive {
				request.Positive = false
				marshal, err := json.Marshal(request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = createJobOrEvent(db, marshal, KeyPEvents, title)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				marshal, err := json.Marshal(request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = createJobOrEvent(db, marshal, KeyNEvents, title)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			lgr.Printf("Event created for %s ", request.Description)

		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode("success")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func addAdminHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request UserInfo
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new admin request: %v", request)

		if request.Email == "" {
			err = fmt.Errorf("Email is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if request.FirstName == "" {
			err = fmt.Errorf("FirstName is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if request.LastName == "" {
			err = fmt.Errorf("LastName is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if request.SchoolId == "" {
			err = fmt.Errorf("SchoolId is mandatory")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		request.Role = UserRoleAdmin
		request.Name = request.Email
		password := RandomString(8)
		request.PasswordSha = EncodePassword(password)

		_, err = CreateSchoolAdmin(db, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		adminTeacher := UserInfo{
			Name:        request.Name[:1] + "." + request.Name[1:],
			FirstName:   request.FirstName,
			LastName:    request.LastName,
			Email:       request.Email[:1] + "." + request.Email[1:],
			Role:        UserRoleTeacher,
			SchoolId:    request.SchoolId,
			PasswordSha: request.PasswordSha,
		}

		err = createTeacher(db, adminTeacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := openapi.ResponseResetPassword{
			Password: password,
		}

		lgr.Printf("admin created for %s ", request.Email)

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func seedDbHandler(db *bolt.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lgr.Printf("INFO new seedDb request")

		err := seedDb(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("seedDB successful")

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode("")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func createJobOrEvent(db *bolt.DB, marshal []byte, bucketKey, itemKey string) error {
	return db.Update(func(tx *bolt.Tx) error {
		EJ, err := tx.CreateBucketIfNotExists([]byte(bucketKey))
		if err != nil {
			return err
		}
		err = EJ.Put([]byte(itemKey), marshal)
		return err
	})
}

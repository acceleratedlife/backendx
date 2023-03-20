package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
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

type ResetPasswordRequest struct {
	Email string
}

type EventRequest struct {
	Positive    bool `json:",omitempty"`
	Description string
	Title       string `json:",omitempty"`
}

func newSchoolHandler(db *bolt.DB, clock Clock) http.Handler {
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
		err = createNewSchool(db, clock, request, response.AdminPassword)
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
		var request ResetPasswordRequest
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

		response, err := resetPassword(db, user, 1)

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

		err = createJobs(db, requests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		var requests []EventRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requests)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new event request: %v", requests)

		err = createEvents(db, requests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode("success")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func seedDbHandler(db *bolt.DB, clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new seedDb request")

		var eventRequests []EventRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&eventRequests)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jobRequests []Job
		err = decoder.Decode(&jobRequests)
		if err != nil {
			err = fmt.Errorf("cannot parse request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = seedDb(db, clock, eventRequests, jobRequests)
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
			Settings: TeacherSettings{
				CurrencyLock: false,
			},
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

func nextCollegeHandler(clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO College request")

		clock.TickOne(time.Hour * 24 * 15)

		lgr.Printf(clock.Now().String())

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode("")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func nextCareerHandler(clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO Career request")

		clock.TickOne(time.Hour * 24 * 5)

		lgr.Printf(clock.Now().String())

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode("")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func nextDayHandler(clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new Day request")

		clock.TickOne(time.Hour * 24)

		lgr.Printf(clock.Now().String())

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode("")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func nextHourHandler(clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new Hour request")

		clock.TickOne(time.Hour)

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode("")
		if err != nil {
			lgr.Printf("ERROR failed to send")
		}

	})
}

func nextMinutesHandler(clock *DemoClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if clock == nil {
			err := fmt.Errorf("This endpoint does not work on the production server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lgr.Printf("INFO new Minute request")

		clock.TickOne(time.Minute * 10)

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode("")
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

func createEvents(db *bolt.DB, events []EventRequest) (err error) {
	for i, request := range events {
		if request.Description == "" {
			return fmt.Errorf("description is mandatory on: " + strconv.Itoa(i))
		}

		if request.Title == "" {
			return fmt.Errorf("title is mandatory on: " + strconv.Itoa(i))
		}

		title := request.Title
		request.Title = ""

		if request.Positive {
			request.Positive = false
			marshal, err := json.Marshal(request)
			if err != nil {
				return err
			}
			err = createJobOrEvent(db, marshal, KeyPEvents, title)
			if err != nil {
				return err
			}
		} else {
			marshal, err := json.Marshal(request)
			if err != nil {
				return err
			}
			err = createJobOrEvent(db, marshal, KeyNEvents, title)
			if err != nil {
				return err
			}
		}

		lgr.Printf("Event created for %s ", request.Description)

	}

	return nil

}

func createJobs(db *bolt.DB, jobs []Job) (err error) {
	for i, request := range jobs {

		if request.Title == "" {
			return fmt.Errorf("title is mandatory on: " + strconv.Itoa(i))
		}

		if request.Pay == 0 {
			return fmt.Errorf("pay is mandatory on: " + strconv.Itoa(i))
		}

		if request.Description == "" {
			return fmt.Errorf("description is mandatory on: " + strconv.Itoa(i))
		}

		title := request.Title
		request.Title = ""

		if request.College {
			request.College = false
			marshal, err := json.Marshal(request)
			if err != nil {
				return err
			}
			err = createJobOrEvent(db, marshal, KeyCollegeJobs, title)
			if err != nil {
				return err
			}
		} else {
			marshal, err := json.Marshal(request)
			if err != nil {
				return err
			}
			err = createJobOrEvent(db, marshal, KeyJobs, title)
			if err != nil {
				return err
			}
		}

		lgr.Printf("job created for %s ", request.Title)

	}

	return nil

}

func createTeachers(db *bolt.DB, clock Clock, schoolId, password string) (err error) {
	for i := 0; i < 30; i++ {
		newTeacher := UserInfo{
			Name:        "tt" + strconv.Itoa(i) + "@tt.com",
			FirstName:   "tt" + strconv.Itoa(i),
			LastName:    "tt" + strconv.Itoa(i),
			Email:       "tt" + strconv.Itoa(i) + "@tt.com",
			Confirmed:   true,
			PasswordSha: EncodePassword(password),
			SchoolId:    schoolId,
			Role:        UserRoleTeacher,
			Settings: TeacherSettings{
				CurrencyLock: false,
			},
		}

		err = createTeacher(db, newTeacher)
		if err != nil {
			return err
		}

		lgr.Printf("Created teacher: %s ", newTeacher.Email)

		for j := 0; j < 2; j++ {
			classId, _, err := CreateClass(db, clock, schoolId, newTeacher.Email, "math"+strconv.Itoa(j), j)
			if err != nil {
				return err
			}

			path := PathId{
				schoolId:  schoolId,
				classId:   classId,
				teacherId: newTeacher.Email,
			}

			for k := 0; k < 5; k++ {
				newStudent := UserInfo{
					Name:        "@ss" + strconv.Itoa((i*30)+k),
					FirstName:   "s" + strconv.Itoa((i*30)+k),
					LastName:    "ss" + strconv.Itoa((i*30)+k),
					Email:       "@ss" + strconv.Itoa((i*30)+k),
					Confirmed:   true,
					PasswordSha: EncodePassword(password),
					SchoolId:    schoolId,
					Role:        UserRoleStudent,
					Job:         getJobId(db, KeyJobs),
				}

				err = createStudent(db, newStudent, path)
				if err != nil {
					return err
				}

				err = pay2Student(db, clock, newStudent, decimal.NewFromInt32(int32(i+k+1)), newTeacher.Email, "start off")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

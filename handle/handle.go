package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tusupov/gousers/db"
	"github.com/tusupov/gousers/model"
)

var (
	errAmountEmpty  = errors.New("`amount` must not be empty")
	errUserNotFound = errors.New("User not found")
)

// UserNew handle function `/user`
// method POST
// body example json: {"amount": 100}
func UserNew(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			log.Printf("handle User panic: %v \ndebug stack: %s\n", err, debug.Stack())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	// Get data from body
	var bodyStruct model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bodyStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new user
	usr, err := db.New(bodyStruct.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write user
	write(w, usr)

}

// User handle function `/user/{id:[0-9]+}`
func User(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			log.Printf("handle User panic: %v \ndebug stack: %s\n", err, debug.Stack())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	// Url variables
	vars := mux.Vars(r)

	// Get id from url
	usrId, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user by id
	usr, ok := db.Get(usrId)
	if !ok {
		http.Error(w, errUserNotFound.Error(), http.StatusNotFound)
		return
	}

	// Write user
	write(w, usr)

}

// UserNew handle function `/user/transfer`
// method POST
// body example json: {"from_user": 1, "to_user": 2, "amount": 100}
func UserAmountTransfer(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			log.Printf("handle User panic: %v \ndebug stack: %s\n", err, debug.Stack())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	// Get data from body
	decoder := json.NewDecoder(r.Body)
	var bodyStruct model.Transfer
	err := decoder.Decode(&bodyStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Transfer
	errTransfer := db.Transfer(bodyStruct.FromUser, bodyStruct.ToUser, bodyStruct.Amount)
	if errTransfer != nil {
		http.Error(w, errTransfer.Error(), http.StatusBadRequest)
		return
	}

	// Write success result
	write(w, map[string]string{"success": "ok"})

}

// write result
func write(w http.ResponseWriter, result interface{}) {

	// Set header status
	w.WriteHeader(http.StatusOK)

	// encode to json
	buf, err := json.Marshal(result)
	if err != nil {
		log.Printf("handle write: %v\ndebug stack: %s\n", err, debug.Stack())
		return
	}

	// write
	if _, err := w.Write(buf); err != nil {
		log.Printf("handle write: %v\ndebug stack: %s\n", err, debug.Stack())
	}

}

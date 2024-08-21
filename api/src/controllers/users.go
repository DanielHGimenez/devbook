package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/security"
	"log"
	"net/http"
	"strconv"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := GetBody(r, &user); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	if err := user.Validate(); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	securedPassword, err := security.Hash(user.Password)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}
	user.Password = string(securedPassword)

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	id, err := repository.Create(user)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusCreated, map[string]uint64{"id": id}, w)
}

func FindAllUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	nick := queryParams.Get("nick")
	email := queryParams.Get("email")

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	users, err := userRepository.FindBy(name, nick, email)
	if err != nil {
		Respond(http.StatusInternalServerError, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	Respond(http.StatusOK, users, w)
}

func FindOneUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	user, err := userRepository.FindOne(userID)
	if err != nil {
		Respond(http.StatusInternalServerError, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	if user != nil {
		user.Password = ""
		Respond(http.StatusOK, user, w)
	} else {
		Respond(http.StatusNotFound, nil, w)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil || userID != tokenUserID {
		Respond(http.StatusForbidden, nil, w)
		return
	}

	var user models.User
	if err := GetBody(r, &user); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	if err := user.Validate(); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	securedPassword, err := security.Hash(user.Password)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}
	user.Password = string(securedPassword)

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	err = userRepository.Update(userID, &user)
	if err != nil {
		Respond(http.StatusInternalServerError, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil || userID != tokenUserID {
		Respond(http.StatusForbidden, nil, w)
		return
	}

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	err = userRepository.Delete(userID)
	if err != nil {
		Respond(http.StatusInternalServerError, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil || userID != tokenUserID {
		Respond(http.StatusForbidden, nil, w)
		return
	}

	var passwordChange models.PasswordChange
	if err := GetBody(r, &passwordChange); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	if err = passwordChange.Validate(); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	currentPassword, err := userRepository.FindPassword(userID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	if err = security.Verify(passwordChange.Current, currentPassword); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Current password is invalid"), w)
		log.Println(err)
		return
	}

	hashNewPassword, err := security.Hash(passwordChange.New)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	if err = userRepository.SavePassword(userID, string(hashNewPassword)); err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

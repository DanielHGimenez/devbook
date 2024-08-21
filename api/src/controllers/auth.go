package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/security"
	"log"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	var authentication models.Authentication
	if err := GetBody(r, &authentication); err != nil {
		Respond(http.StatusUnprocessableEntity, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	if err := authentication.Validar(); err != nil {
		Respond(http.StatusBadRequest, err.Error(), w)
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

	users, err := userRepository.FindBy("", "", authentication.Email)
	if err != nil {
		Respond(http.StatusInternalServerError, NewErrorResponse(err.Error()), w)
		log.Println(err)
		return
	}

	if len(users) == 0 {
		Respond(http.StatusNotFound, nil, w)
		log.Println(err)
		return
	}

	userFound := users[0]
	if security.Verify(authentication.Password, userFound.Password) != nil {
		Respond(http.StatusUnauthorized, nil, w)
		log.Println(err)
		return
	}

	jwt, err := security.CreateToken(userFound.ID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusAccepted, jwt, w)
}

package controllers

import (
	"api/src/database"
	"api/src/repositories"
	"api/src/security"
	"log"
	"net/http"
	"strconv"
)

func Follow(w http.ResponseWriter, r *http.Request) {
	followedID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	followerID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}

	if followerID == followedID {
		Respond(http.StatusBadRequest, NewErrorResponse("You can't follow yourself"), w)
		return
	}

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	followerRepository := repositories.NewFollowerRepository(db)

	wasFollow, err := followerRepository.Exists(followedID, followerID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	if !wasFollow {
		err = followerRepository.Create(followedID, followerID)
		if err != nil {
			Respond(http.StatusInternalServerError, nil, w)
			log.Println(err)
			return
		}
	}

	Respond(http.StatusOK, nil, w)
}

func UnFollow(w http.ResponseWriter, r *http.Request) {
	followedID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	followerID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}

	if followerID == followedID {
		Respond(http.StatusBadRequest, NewErrorResponse("You can't unfollow yourself"), w)
		return
	}

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	followerRepository := repositories.NewFollowerRepository(db)

	err = followerRepository.Delete(followedID, followerID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

func FindAllFollowers(w http.ResponseWriter, r *http.Request) {
	followedID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
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

	followers, err := userRepository.FindAllFollowers(followedID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, followers, w)
}

func FindAllFollowing(w http.ResponseWriter, r *http.Request) {
	followerID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
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

	followers, err := userRepository.FindAllFollowing(followerID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, followers, w)
}

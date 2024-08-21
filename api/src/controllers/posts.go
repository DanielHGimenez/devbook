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

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := GetBody(r, &post); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		return
	}
	post.AuthorID = tokenUserID

	db, err := database.Connect()
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	defer db.Close()

	postRepository := repositories.NewPostRepository(db)

	id, err := postRepository.Create(post)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusCreated, map[string]uint64{"id": id}, w)
}

func FindRelatedPosts(w http.ResponseWriter, r *http.Request) {
	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	posts, err := postRepository.FindRelatedPosts(tokenUserID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, posts, w)
}

func FindOnePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
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

	postRepository := repositories.NewPostRepository(db)

	post, err := postRepository.FindOne(postID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	if post == nil {
		Respond(http.StatusNotFound, nil, w)
		return
	}

	Respond(http.StatusOK, post, w)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	var post models.Post
	if err := GetBody(r, &post); err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Body required"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	isAuthor, err := postRepository.IsAuthor(tokenUserID, postID)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}
	if !isAuthor {
		Respond(http.StatusForbidden, nil, w)
		return
	}

	err = postRepository.Update(postID, post)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	isAuthor, err := postRepository.IsAuthor(tokenUserID, postID)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}
	if !isAuthor {
		Respond(http.StatusForbidden, nil, w)
		return
	}

	err = postRepository.Delete(postID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, nil, w)
}

func WasPostLiked(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	exists, err := postRepository.Exists(postID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	if !exists {
		Respond(http.StatusNotFound, nil, w)
		return
	}

	wasLiked, err := postRepository.WasLiked(tokenUserID, postID)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}

	Respond(http.StatusOK, map[string]bool{"wasLiked": wasLiked}, w)
}

func Like(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	exists, err := postRepository.Exists(postID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	if !exists {
		Respond(http.StatusNotFound, nil, w)
		return
	}

	wasLiked, err := postRepository.WasLiked(tokenUserID, postID)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}
	if !wasLiked {
		err = postRepository.Like(tokenUserID, postID)
		if err != nil {
			Respond(http.StatusInternalServerError, nil, w)
			log.Println(err)
			return
		}
	}

	Respond(http.StatusCreated, nil, w)
}

func UnLike(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(GetPathParam("id", r), 10, 32)
	if err != nil {
		Respond(http.StatusBadRequest, NewErrorResponse("Failed to parse path param"), w)
		log.Println(err)
		return
	}

	tokenUserID, err := strconv.ParseUint(r.Header[security.UserIDHeader][0], 10, 32)
	if err != nil {
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

	postRepository := repositories.NewPostRepository(db)

	exists, err := postRepository.Exists(postID)
	if err != nil {
		Respond(http.StatusInternalServerError, nil, w)
		log.Println(err)
		return
	}
	if !exists {
		Respond(http.StatusNotFound, nil, w)
		return
	}

	wasLiked, err := postRepository.WasLiked(tokenUserID, postID)
	if err != nil {
		Respond(http.StatusForbidden, nil, w)
		log.Println(err)
		return
	}
	if wasLiked {
		err = postRepository.UnLike(tokenUserID, postID)
		if err != nil {
			Respond(http.StatusInternalServerError, nil, w)
			log.Println(err)
			return
		}
	}

	Respond(http.StatusOK, nil, w)
}

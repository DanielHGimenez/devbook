package routes

import (
	"api/src/controllers"
	"net/http"
)

var PostsRoutes = []Route{
	{
		URI:                      "/posts",
		Method:                   http.MethodPost,
		Func:                     controllers.CreatePost,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts",
		Method:                   http.MethodGet,
		Func:                     controllers.FindRelatedPosts,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}",
		Method:                   http.MethodGet,
		Func:                     controllers.FindOnePost,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}",
		Method:                   http.MethodPut,
		Func:                     controllers.UpdatePost,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}",
		Method:                   http.MethodDelete,
		Func:                     controllers.DeletePost,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}/like",
		Method:                   http.MethodGet,
		Func:                     controllers.WasPostLiked,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}/like",
		Method:                   http.MethodPost,
		Func:                     controllers.Like,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/posts/{id}/like",
		Method:                   http.MethodDelete,
		Func:                     controllers.UnLike,
		IsAuthenticationRequired: true,
	},
}

package routes

import (
	"api/src/controllers"
	"net/http"
)

var UsersRoutes = []Route{
	{
		URI:                      "/users",
		Method:                   http.MethodPost,
		Func:                     controllers.CreateUser,
		IsAuthenticationRequired: false,
	},
	{
		URI:                      "/users",
		Method:                   http.MethodGet,
		Func:                     controllers.FindAllUsers,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}",
		Method:                   http.MethodGet,
		Func:                     controllers.FindOneUser,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}",
		Method:                   http.MethodPut,
		Func:                     controllers.UpdateUser,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}",
		Method:                   http.MethodDelete,
		Func:                     controllers.DeleteUser,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}/followers",
		Method:                   http.MethodPost,
		Func:                     controllers.Follow,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}/followers",
		Method:                   http.MethodDelete,
		Func:                     controllers.UnFollow,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}/followers",
		Method:                   http.MethodGet,
		Func:                     controllers.FindAllFollowers,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/followers/users/{id}",
		Method:                   http.MethodGet,
		Func:                     controllers.FindAllFollowing,
		IsAuthenticationRequired: true,
	},
	{
		URI:                      "/users/{id}/password",
		Method:                   http.MethodPost,
		Func:                     controllers.ChangePassword,
		IsAuthenticationRequired: true,
	},
}

package routes

import (
	"api/src/controllers"
	"net/http"
)

var AuthRoutes = []Route{
	{
		URI:                      "/signin",
		Method:                   http.MethodPost,
		Func:                     controllers.SignIn,
		IsAuthenticationRequired: false,
	},
}

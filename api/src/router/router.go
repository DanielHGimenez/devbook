package router

import (
	"api/src/router/middlewares"
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

func Build() *mux.Router {
	router := mux.NewRouter()

	routeList := routes.UsersRoutes
	routeList = append(routeList, routes.AuthRoutes...)
	routeList = append(routeList, routes.PostsRoutes...)

	for _, route := range routeList {
		if route.IsAuthenticationRequired {
			router.HandleFunc(route.URI, middlewares.Authorize(route.Func)).Methods(route.Method)
		} else {
			router.HandleFunc(route.URI, route.Func).Methods(route.Method)
		}
	}
	return router
}

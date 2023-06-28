package api

import (
	"github.com/jeremyauchter/adjutor/api/controllers"
	"github.com/jeremyauchter/adjutor/api/routes"
)

var server = controllers.Server{}
var router = routes.Routers{}

func Run() {
	server.Initialize()
	router.StartRouter()
	router.InitializeRoutes(server)
	router.InitializeTagRoutes(server)
	router.InitializeVendorRoutes(server)
	router.Run(":8080")
}

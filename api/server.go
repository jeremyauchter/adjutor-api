package api

import (
	"github.com/jeremyauchter/adjutor/api/controllers"
	"github.com/jeremyauchter/adjutor/api/routes"
)

var server = controllers.Server{}
var router = routes.Routers{}

func startRoutes() {
	router.InitializeRoutes(server)
	router.InitializeTagRoutes(server)
	router.InitializeVendorRoutes(server)
	router.InitializeAudienceRoutes(server)
	router.InitializeDepartmentRoutes(server)
	router.InitializeCountryRoutes(server)
	router.InitializeItemVariantRoutes(server)
	router.InitializeClassRoutes(server)
	router.InitializeProductRoutes(server)
	router.InitializeProductTypeRoutes(server)
	router.InitializeStyleRoutes(server)
}

func Run() {
	server.Initialize()
	router.StartRouter()
	startRoutes()
	router.Run(":8080")
}

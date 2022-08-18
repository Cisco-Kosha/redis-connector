package app

import (
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) initializeRoutes() {
	var apiV1 = "/api/v1"

	// specification routes
	a.Router.HandleFunc(apiV1+"/specification/list", a.listConnectorSpecification).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/specification/test", a.testConnectorSpecification).Methods("POST", "OPTIONS")

	a.Router.HandleFunc(apiV1+"/keys", a.getAllKeys).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/ping", a.ping).Methods("GET", "OPTIONS")
	// kv
	a.Router.HandleFunc(apiV1+"/kv/{key}", a.getSingleKey).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/kv/{key}", a.setValueForSingleKey).Methods("POST", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/kv/{key}", a.removeSingleKey).Methods("DELETE", "OPTIONS")

	// set
	a.Router.HandleFunc(apiV1+"/set/{key}", a.getSetElements).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/set/{key}", a.addElementsToASet).Methods("POST", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/set/{key}", a.removeElementsFromASet).Methods("DELETE", "OPTIONS")

	// list
	a.Router.HandleFunc(apiV1+"/list/{key}", a.getListElements).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/list/{key}", a.addElementsToAList).Methods("POST", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/list/{key}", a.removeElementsFromAList).Methods("DELETE", "OPTIONS")

	// hashes
	a.Router.HandleFunc(apiV1+"/hash/{key}", a.getHashElements).Methods("GET", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/hash/{key}", a.addElementsToAHash).Methods("POST", "OPTIONS")
	a.Router.HandleFunc(apiV1+"/hash/{key}", a.removeElementsFromAHash).Methods("DELETE", "OPTIONS")

	// Swagger
	a.Router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
}

package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	header := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//handlers.CORS(header, methods, origins)(router)

	router.HandleFunc("/apiV1/saveLoraData", saveLoraData)
	router.HandleFunc("/apiV1/devices/show", showDevices).Methods("GET")
	router.HandleFunc("/apiV1/devices/show/{id}", showDevice).Methods("GET")
	router.HandleFunc("/apiV1/devices/add", addDevice).Methods("POST")
	router.HandleFunc("/apiV1/devices/delete/{id}", deleteDevice).Methods("DELETE")
	router.HandleFunc("/downlink", downlink)
	router.HandleFunc("/apiV1/devices/update/{id}", updateDevice).Methods("PUT")

	readConfig()
	printConfig()
	//createOutputFiles()
	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(config.Listeningport, 10), handlers.CORS(header, methods, origins)(router)))
	//log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(config.Listeningport, 10), router))
	//fpOnYieldData.Close()
}

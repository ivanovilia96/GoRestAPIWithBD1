package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	defer ConnectedDataBase.Close()
	println("You can work on http://localhost:8080/ ")
	r := mux.NewRouter()
	r.HandleFunc("/notifications/page={pageNumber}", GetNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notifications/page={pageNumber}/sort/price={priceSortType}/date={dateSortType}", GetNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}", GetNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}/optionalFields={fields}", GetNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification", PutNotification).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", r))
}

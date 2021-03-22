package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	defer ConnectedDataBase.Close()
	println("You can work on http://localhost:8081/ ")
	r := mux.NewRouter()
	r.HandleFunc("/notifications/page={pageNumber}", GetNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notifications/page={pageNumber}/sort/price={priceSortType}/date={dateSortType}", GetNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}", GetNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}/optionalFields={fields}", GetNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification", PutNotification).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", r))
}

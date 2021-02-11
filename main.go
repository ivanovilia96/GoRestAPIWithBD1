package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	defer ConnectedDataBase.Close()
	r := mux.NewRouter()
	r.HandleFunc("/notifications/page={pageNumber}", getNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notifications/page={pageNumber}/sort/price={priceSortType}/date={dateSortType}", getNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}", getNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification/{id}/optionalFields={fields}", getNotification).Methods(http.MethodGet)
	r.HandleFunc("/notification", putNotification).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", r))
}

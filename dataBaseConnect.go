package main

import "database/sql"

func dataBaseConnect() *sql.DB {
	var (
		login            = "root"
		password         = "root"
		connectionMethod = "@tcp"
		hostname         = "127.0.0.1"
		port             = "3306"
		DBName           = "firstDB"
	)
	db, err := sql.Open("mysql", login+":"+password+connectionMethod+"("+hostname+":"+port+")/"+DBName)
	if err != nil {
		panic(err)
	}
	return db
}

var ConnectedDataBase = dataBaseConnect()

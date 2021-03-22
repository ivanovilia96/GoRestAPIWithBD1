package main

import (
	"database/sql"
)

func dataBaseConnect() *sql.DB {
	var (
		connectionMethod = "@tcp"
		hostname         = "host.docker.internal"
		port             = "3306"
		DBName           = "firstDB"
		login, password  = "root", "root"
	)
	db, err := sql.Open("mysql", login+":"+password+connectionMethod+"("+hostname+":"+port+")/")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DBName)
	if err != nil {
		panic("Error on creating dataBase:" + err.Error())
	}

	_, err = db.Exec("USE " + DBName)
	if err != nil {
		panic("Error when we call USE dataBase : " + err.Error())
	}

	sqlQueryCreateNotes := "CREATE TABLE IF NOT EXISTS Notes(ID INT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE,price int NOT NULL,placementData DATE NOT NULL,name varchar(200) not null,description text,PRIMARY KEY(ID) )"
	sqlQueryCreateImagesForNotes := "CREATE TABLE IF NOT EXISTS ImagesForNotes( image_id INT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE,note_id INT UNSIGNED NOT NULL,image_data text,CONSTRAINT any PRIMARY KEY(image_id), FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE)"
	_, err = db.Exec(sqlQueryCreateNotes)
	if err != nil {
		panic("Error on creating Notes table: " + err.Error())
	}
	_, err = db.Exec(sqlQueryCreateImagesForNotes)
	if err != nil {
		panic("Error on creating ImagesForNotes table: " + err.Error())
	}

	return db
}

var ConnectedDataBase = dataBaseConnect()

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var (
		countNotesOnOnePage int64 = 10
		from                int64 = 0
		to                  int64 = 10
		priceSortType             = ""
		dateSortType              = ""
		vars                      = mux.Vars(r)
		resultData          []Note
		pageNumber, err     = strconv.ParseInt(vars["pageNumber"], 10, 32)
	)

	if err != nil {
		panic(err.Error())
	}

	if pageNumber > 1 {
		from = (countNotesOnOnePage * pageNumber) - countNotesOnOnePage
		to = countNotesOnOnePage * pageNumber
	}

	if vars["priceSortType"] == "max-to-min" {
		priceSortType = " desc "
	}

	if vars["dateSortType"] == "max-to-min" {
		dateSortType = " desc "
	}

	var sqlStatement string
	sqlStatement = "SELECT  n.name, n.price, im.image_data, n.placementdata, n.id FROM notes n, imagesfornotes im " +
		"where im.note_id = n.id group by id, placementdata order by n.price" + priceSortType +
		", n.placementdata " + dateSortType + " limit " + strconv.FormatInt(from, 10) + `,` + strconv.FormatInt(to, 10)

	results, err := ConnectedDataBase.Query(sqlStatement)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var (
			oneRowDB Note
		)
		oneRowDB.ImageData = append(oneRowDB.ImageData, "")
		err = results.Scan(&oneRowDB.Name, &oneRowDB.Price, &oneRowDB.ImageData[0], &oneRowDB.PlacementData, &oneRowDB.Id)
		if err != nil {
			panic(err.Error())
		}
		resultData = append(resultData, oneRowDB)
	}

	resultDataJson, err := json.Marshal(
		struct {
			Data []getNoteResp
		}{
			Data: getFieldsFromArray(resultData),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultDataJson)
}

func getNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var (
		vars                 = mux.Vars(r)
		optionalFieldsForSql = ""
		needAllImages        = false
		optionalFields       = strings.Split(vars["fields"], ",")
	)

	for _, v := range optionalFields {
		if v == "description" {
			optionalFieldsForSql += ", n.description "
		}
		if v == "allImages" {
			needAllImages = true
		}
	}

	var (
		oneNote      getNoteResp
		sqlStatement = "SELECT  n.name, n.price " + optionalFieldsForSql + ", im.image_data FROM notes n, imagesfornotes im " +
			"where im.note_id = n.id and n.id = " + vars["id"] + " group by id "
	)

	if needAllImages {
		sqlStatementForAllImages := "select ImagesForNotes.image_data from ImagesForNotes where ImagesForNotes.note_id = " + vars["id"]
		if optionalFieldsForSql == "" {
			err := ConnectedDataBase.QueryRow(sqlStatement).Scan(&oneNote.Name, &oneNote.Price, &oneNote.ImageData)
			if err != nil {
				err.Error()
			}
		} else {
			err := ConnectedDataBase.QueryRow(sqlStatement).Scan(&oneNote.Name, &oneNote.Price, &oneNote.Description, &oneNote.ImageData)
			if err != nil {
				err.Error()
			}
		}

		results, err := ConnectedDataBase.Query(sqlStatementForAllImages)
		if err != nil {
			panic(err.Error())
		}
		var localStoreData []string
		for results.Next() {
			var localStr string
			err = results.Scan(&localStr)
			if err != nil {
				panic(err.Error())
			}
			localStoreData = append(localStoreData, localStr)

		}
		oneNote.ImageData = localStoreData
	} else {
		oneNote.ImageData = append(oneNote.ImageData, "")
		if vars["fields"] != "" {
			err := ConnectedDataBase.QueryRow(sqlStatement).Scan(&oneNote.Name, &oneNote.Price, &oneNote.Description, &oneNote.ImageData[0])
			if err != nil {
				err.Error()
			}
		} else {
			err := ConnectedDataBase.QueryRow(sqlStatement).Scan(&oneNote.Name, &oneNote.Price, &oneNote.ImageData[0])
			if err != nil {
				err.Error()
			}
		}

	}

	resultDataJson, err := json.Marshal(
		struct {
			Data getNoteResp
		}{
			Data: oneNote,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultDataJson)
}

type NotificationPut struct {
	Name        string
	Description string
	Image_data  []string
	Price       int
}

func putNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	bodyInStructure := NotificationPut{}
	err = json.Unmarshal(requestBody, &bodyInStructure)

	if err != nil {
		panic(err)
	}
	if len(bodyInStructure.Image_data) > 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Нельзя помещать более 3х ссылок на обьявление"}`))
	} else if len(bodyInStructure.Description) > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Размер описания не должен быть более 1000 символов"}`))
	} else if len(bodyInStructure.Name) > 200 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Наименование обьявления не должно быть более 200 символов"}`))
	} else {
		sqlQuery := "insert into notes(price,name,placementdata,description) values (\"" + strconv.Itoa(bodyInStructure.Price) + "\" ," +
			"\"" + bodyInStructure.Name + "\", \"" + (time.Now().Format("2006-01-02")) + "\",\"" + bodyInStructure.Description + "\")"

		res, err := ConnectedDataBase.Exec(sqlQuery)
		if err != nil {
			panic(err.Error())
		}

		id, err := res.LastInsertId()
		if err != nil {
			panic(err.Error())
		}
		queryForImages := "insert into ImagesForNotes(note_id, image_data) values "
		for i, v := range bodyInStructure.Image_data {
			queryForImages += "( \"" + strconv.Itoa(int(id)) + "\", \" " + v + "\" )"

			if len(bodyInStructure.Image_data) == i+1 {
				queryForImages += ";"
			} else {
				queryForImages += ","
			}
		}
		_, err = ConnectedDataBase.Exec(queryForImages)
		if err != nil {
			err.Error()
		}
		dataJson, err := json.Marshal(struct {
			Id     string
			Status int
		}{strconv.Itoa(int(id)), http.StatusAccepted})
		if err != nil {
			err.Error()
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(dataJson)
	}
}

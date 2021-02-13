package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type RequestData struct {
	Data []getNoteResp
}

var (
	uniqueName        = "QSA|?FSA/.L:A||sf]}"
	onWhatIdWeHaveRow = "67"
)

func TestGetNotificationsWithoutSort(t *testing.T) {
	req, err := http.NewRequest("GET", "/notifications/page=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notifications/page={pageNumber}", GetNotifications)
	router.ServeHTTP(recorder, req)

	// -> МОЖНО ПРОВЕРИТЬ ОТВЕТ, ( внизу пример как получить данные и структура ) , но,
	// тогда придется сначала теста вызывать функцию добавления в бд данных Х раз (что бы данные возвратились, а то если БД пустая, то ошибка неожиданная будет)
	//+ как то искать именно эти данные которые мы добавили, ведь в данном запросе допустим, возвращается список, то есть по нему итерироваться нужно и тд,
	//а после теста удаления данных из бд.
	//возможно  есть какой то более простой способ протестировать тело ответа ( может тестовые бд онлайн какие нибудь существуют)

	if recorder.Code != http.StatusOK {
		t.Error("should work")
	}

}

//сделал тут через имя потому что нет доступа к LastInsertId т.к. он не возвращается из PutNotification fun
func deleteRowAfterTest() {
	_, err := ConnectedDataBase.Exec("DELETE FROM notes WHERE name=\"" + uniqueName + "\";")
	if err != nil {
		panic(err.Error())
	}
}

func createRowInTableForTest(t *testing.T) {
	group := getNoteResp{
		uniqueName,
		[]sql.NullString{},
		19,
		"desc",
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.NewRequest("PUT", "/notification", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PutNotification)

	handler.ServeHTTP(rr, res)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler PutNotification returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "\"Status\":202") {
		t.Errorf("handler returned unexpected body: got %v want containing string \"Status\":202 ",
			rr.Body.String())
	}
}

func TestPutNotification(t *testing.T) {
	createRowInTableForTest(t)
	deleteRowAfterTest()

}

func TestGetNotificationsWithSort(t *testing.T) {
	// добавляем колонку для теста c помощью фактически вызова другой функции ( она уже вызывается до, да, но может переместиться )
	// это сделно для проверки, что возвращается хотя бы запись, (мы её ниже добавляем)
	createRowInTableForTest(t)
	defer deleteRowAfterTest()

	// наш тест
	req, err := http.NewRequest("GET", "/notifications/page=1/sort/price=max-to-min/date=max-to-min", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notifications/page={pageNumber}/sort/price={priceSortType}/date={dateSortType}", GetNotifications)
	router.ServeHTTP(recorder, req)
	requestBody, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		panic(err)
	}
	bodyInStructure := RequestData{}
	err = json.Unmarshal(requestBody, &bodyInStructure)

	if len(bodyInStructure.Data) == 0 {
		t.Errorf("data isn`t return")
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("response code is incorrect: %v", recorder.Code)
	}
}

func TestGetNotificationWithoutOptionalFields(t *testing.T) {
	req, err := http.NewRequest("GET", "/notification/"+onWhatIdWeHaveRow, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notification/{id}", GetNotification)
	router.ServeHTTP(recorder, req)
	requestBody, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		panic(err)
	}
	bodyInStructure := struct {
		Data getNoteResp
	}{}
	err = json.Unmarshal(requestBody, &bodyInStructure)
	if bodyInStructure.Data.Name == "" || bodyInStructure.Data.Price == 0 {
		t.Errorf("data isn`t return, ВОЗМОЖНО ВЫ ЗАБЫЛИ ПОМЕНЯТЬ ПЕРЕМЕННУЮ onWhatIdWeHaveRow на ваш id")
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("response code is incorrect: %v", recorder.Code)
	}

}

func TestGetNotificationWithOptionalFields(t *testing.T) {
	req, err := http.NewRequest("GET", "/notification/"+onWhatIdWeHaveRow+"/optionalFields=description,allImages", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notification/{id}/optionalFields={fields}", GetNotification)
	router.ServeHTTP(recorder, req)
	requestBody, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		panic(err)
	}
	bodyInStructure := struct {
		Data getNoteResp
	}{}
	err = json.Unmarshal(requestBody, &bodyInStructure)

	if bodyInStructure.Data.Name == "" || bodyInStructure.Data.Price == 0 {
		t.Errorf("data isn`t return, ВОЗМОЖНО ВЫ ЗАБЫЛИ ПОМЕНЯТЬ ПЕРЕМЕННУЮ onWhatIdWeHaveRow на ваш id")
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("response code is incorrect: %v", recorder.Code)
	}
}

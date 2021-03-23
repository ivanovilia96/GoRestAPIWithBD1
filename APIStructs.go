package main

import (
	"database/sql"
	"encoding/json"
)

type NullString sql.NullString

func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

type NotificationPut struct {
	Name        string
	Description string
	Image_data  []string
	Price       int
}

type getNoteResp struct {
	Name        string       `json:"name"`
	ImageData   []NullString `json:"imageData"`
	Price       int          `json:"price"`
	Description string       `json:"description,omitempty"`
}

// решение использовать данный тип NullString вызвано тем, что если пользователь не добавит не 1 картинки, то  будет выброшена ошибка связанная с null возвращаемым типо данных
// - возвращать и написать что картинок т.е в скане с sql запроса возвращался null, потому что картинок нет у пользователя, а null нельзя к string
// на основании этого, воспользовался гуглом и нашел выход, ( https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267 )
//возможно есть какой то более лаконичный выход

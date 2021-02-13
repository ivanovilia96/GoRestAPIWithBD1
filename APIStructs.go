package main

import "database/sql"

type NotificationPut struct {
	Name        string
	Description string
	Image_data  []string
	Price       int
}

type getNoteResp struct {
	Name        string
	ImageData   []sql.NullString
	Price       int
	Description string
}

// решение использовать данный тип NullString вызвано тем, что если пользователь не добавит не 1 картинки, то  выход может быть такой

// - возвращать и написать что картинок т.е в скане с sql запроса возвращался null, потому что картинок нет у пользователя, а null нельзя к string
// на основании этого, воспользовался гуглом и нашел выход, ( https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267 )
//возможно есть какой то более лаконичный выход

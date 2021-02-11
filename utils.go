package main

// Используется для того, что бы выбрать из нашей структуры Note поля , которые мы возвратим в ответе, тут можно их переименовать их и тд (возвратить структуру другого обьекта)
func getFieldsFromArray(arrNotes []Note) []getNoteResp {
	var result []getNoteResp

	for _, v := range arrNotes {
		result = append(result, getNoteResp{
			v.Name,
			v.ImageData,
			v.Price,
			"",
		})
	}
	return result
}

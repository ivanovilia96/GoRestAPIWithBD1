package main

type getNoteResp struct {
	Name        string
	ImageData   []string
	Price       int
	Description string
}

type Note struct {
	Id            int
	Price         int
	PlacementData string
	Name          string
	ImageData     []string
}

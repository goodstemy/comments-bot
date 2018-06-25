package main

type bot struct {
	AccessToken      string `json:"access_token"`
	GroupID          int
	ResponseWall     response
	ResponseComments response
}

type response struct {
	Wall wall `json:"response"`
}

type wall struct {
	Posts []post `json:"items"`
}

type post struct {
	ID int `json:"id"`
}

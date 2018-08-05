package main

// Bot is a main struct for parsing
type Bot struct {
	Config 		*Config
	GroupList 	[]string
	Resp        map[int]*Response
	Result      []ResultPost
}

type ResultPost struct {
	Body    string
	Comment string
}

type Response struct {
	ResponseWall     responseWall
	ResponseComments responseComments
}

type responseWall struct {
	Wall wall `json:"response"`
}

type wall struct {
	Posts []post `json:"items"`
}

type post struct {
	ID           int `json:"id"`
	TopCommentID int
	Likes        like   `json:"likes,omitempty"`
	Text         string `json:"text,omitempty"`
}

type responseComments struct {
	CommentsList comments `json:"response"`
}

type comments struct {
	Comments []comment `json:"items"`
}

type comment struct {
	ID    int    `json:"id"`
	Likes like   `json:"likes,omitempty"`
	Text  string `json:"text,omitempty"`
}

type like struct {
	Count int `json:"count"`
}

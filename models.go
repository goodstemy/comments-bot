package main

import "github.com/Syfaro/telegram-bot-api"

// Bot is a main struct for parsing
type Bot struct {
	Config      *Config
	GroupList   []string
	BotInstance *tgbotapi.BotAPI
}

type ResponseGroupId struct {
	Response []GroupId `json:"response"`
}

type GroupId struct {
	ID int `json:"id"`
}

type ResponseWall struct {
	Response ResponseWallItems `json:"response"`
}

type ResponseWallItems struct {
	Items []ResponseWallId `json:"items"`
}

type ResponseWallId struct {
	ID int `json:"id"`
}

type ResponseComments struct {
	Response ResponseCommentsItems `json:"response"`
}

type ResponseCommentsItems struct {
	Items []ResponseCommentId `json:"items"`
}

type ResponseCommentId struct {
	ID    int                   `json:"id"`
	Likes ResponseCommentsLikes `json:"likes"`
}

type ResponseCommentsLikes struct {
	Count int `json:"count"`
}

type ResultGroupData struct {
	GroupId int
	Posts   []PostWithBestComment
}

type PostWithBestComment struct {
	PostID    int
	CommentID int
}

type ResultContentOfPost struct {
	Response []ResponseContentOfPost `json:"response"`
}

type ResponseContentOfPost struct {
	Text        string             `json:"text"`
	Attachments []AttachmentOfPost `json:"attachments"`
}

type AttachmentOfPost struct {
	Type  string                    `json:"type"`
	Photo AttachmentOfPostPhotoType `json:"photo,omitempty"`
}

type AttachmentOfPostPhotoType struct {
	ID      int `json:"id"`
	OwnerID int `json:"owner_id"`
}

type ToPost struct {
	GroupID map[int][]GroupContentToPost
}

type GroupContentToPost struct {
	Text        string
	CommentText string
	Photos      []string
}

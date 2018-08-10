package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"errors"
	"time"
	"github.com/Syfaro/telegram-bot-api"
)

type Config struct {
	VKAccessToken     string
	TGAccessToken     string
	GroupListFileName string
	TGChatID          int64
	Port              int
	IsProduction      bool
}

// BaseAPIURL is base url for all requests
const BaseAPIURL = "https://api.vk.com/method/"

func main() {
	cfg := Config{}
	err := getConfig(&cfg)

	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TGAccessToken)

	if err != nil {
		panic(err)
	}

	b := Bot{
		Config:      &cfg,
		BotInstance: bot,
	}

	b.BotInstance.Debug = true

	err = b.getGroupList()

	if err != nil {
		panic(err)
	}

	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//updates, err := b.BotInstance.GetUpdatesChan(u)
	//
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	b.BotInstance.Send(msg)
	//}

	var resultGroupsData []ResultGroupData
	start := time.Now()
	parsingComplete := make(chan bool, len(b.GroupList))

	for _, groupName := range b.GroupList {
		go func() {
			r := ResultGroupData{}

			groupId, err := getGroupId(groupName, b.Config.VKAccessToken)

			responsePostsIdList, err := b.getPostsByGroupName(groupId)

			if err != nil {
				// TODO: error handling
				log.Println(err)
			}

			postsIdList := responsePostsIdList.Response.Items;
			for _, post := range postsIdList {
				commentId, err := b.getBestCommentIdOfPost(groupId, post.ID)

				if err != nil {
					// TODO: error handling
					log.Println(err)
					return
				}

				r.Posts = append(r.Posts, PostWithBestComment{
					PostID:    post.ID,
					CommentID: commentId,
				})

			}

			r.GroupId = groupId
			resultGroupsData = append(resultGroupsData, r)

			parsingComplete <- true
		}()
	}

	elapsed := time.Since(start)
	channelsCounter := 0

	for <-parsingComplete {
		channelsCounter++

		if channelsCounter == len(b.GroupList) {
			close(parsingComplete)
		}
	}

	log.Printf("Parsing of %d groups took %s", len(b.GroupList), elapsed)

	for _, group := range resultGroupsData {
		for _, post := range group.Posts {
			b.getContentOfPost(group.GroupId, post.PostID)
		}
	}
}

func (b *Bot) getGroupList() error {
	dir, _ := os.Getwd()
	file, err := os.Open(fmt.Sprintf("%s/%s", dir, b.Config.GroupListFileName))
	defer file.Close()

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		b.GroupList = append(b.GroupList, scanner.Text())
	}

	return nil
}

func (b *Bot) getPostsByGroupName(groupId int) (ResponseWall, error) {
	r := ResponseWall{}

	if groupId == 0 {
		return r, errors.New("Group id cannot be equal 0")
	}

	getPostsByGroupIDURL := fmt.Sprintf(
		"%swall.get?owner_id=%d&count=10&v=5.52&access_token=%s",
		BaseAPIURL, groupId, b.Config.VKAccessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return r, err
	}

	json.Unmarshal(body, &r)

	return r, nil
}

func (b *Bot) getBestCommentIdOfPost(groupId int, postId int) (int, error) {
	getCommentsByIDURL := fmt.Sprintf(
		"%swall.getComments?owner_id=%d&post_id=%d&"+
			"need_likes=1&count=100&v=5.52&access_token=%s",
		BaseAPIURL, groupId, postId, b.Config.VKAccessToken)

	body, err := executeRequest(getCommentsByIDURL)

	if err != nil {
		return 0, err
	}

	r := ResponseComments{}

	json.Unmarshal(body, &r)

	bestLikes := 0
	commentId := 0

	for _, comment := range r.Response.Items {
		if comment.Likes.Count > bestLikes {
			commentId = comment.ID
			bestLikes = comment.Likes.Count
		}
	}

	return commentId, nil
}

func (b *Bot) getContentOfPost(groupId int, postId int) {
	getPostByGroupURL := fmt.Sprintf(
		"%swall.getById?posts=%d_%d&v=5.52&access_token=%s",
		BaseAPIURL, groupId, postId, b.Config.VKAccessToken)

	body, err := executeRequest(getPostByGroupURL)

	if err != nil {
		log.Println(err)
	}

	r := ResultContentOfPost{}

	json.Unmarshal(body, &r)

	log.Println("%+v", r)

	for _, resp := range r.Response {
		msg := tgbotapi.NewMessage(b.Config.TGChatID, resp.Text)
		b.BotInstance.Send(msg);
	}

}

func getGroupId(groupName string, accessToken string) (int, error) {
	getPostsByGroupIDURL := fmt.Sprintf(
		"%sgroups.getById?group_id=%s&v=5.52&access_token=%s",
		BaseAPIURL, groupName, accessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return 0, err
	}

	r := ResponseGroupId{}

	json.Unmarshal(body, &r)

	return -r.Response[0].ID, nil
}

func getConfig(cfg *Config) error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func executeRequest(URL string) ([]byte, error) {
	resp, err := http.Get(URL)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, err
}

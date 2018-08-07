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
)

type Config struct {
	AccessToken       string
	GroupListFileName string
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

	b := Bot{
		Config: &cfg,
	}

	err = b.getGroupList()

	if err != nil {
		panic(err)
	}

	for _, groupName := range b.GroupList {
		groupId, err := getGroupId(groupName, b.Config.AccessToken)

		responsePostsIdList, err := b.getPostsByGroupName(groupId)

		if err != nil {
			// TODO: error handling
			log.Println(err)
		}

		postsIdList := responsePostsIdList.Response.Items;

		for _, postId := range postsIdList {
			commentId, err := b.getBestCommentIdOfPost(groupId, postId.ID)

			if err != nil {
				// TODO: error handling
				log.Println(err)
			}

			log.Println(postId.ID, commentId)
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
		BaseAPIURL, groupId, b.Config.AccessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return r, err
	}

	json.Unmarshal(body, &r)

	return r, nil
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

func getGroupId(groupName string, accessToken string) (int, error) {
	log.Println(accessToken)
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

func (b *Bot) getBestCommentIdOfPost(groupId int, postId int) (int, error) {
	getCommentsByIDURL := fmt.Sprintf(
		"%swall.getComments?owner_id=%d&post_id=%d&"+
			"need_likes=1&count=100&v=5.52&access_token=%s",
		BaseAPIURL, groupId, postId, b.Config.AccessToken)

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

package main

import (
	"os"
	"bufio"
		"fmt"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	)

type Config struct {
	AccessToken  	  string
	GroupListFileName string
	Port         	  int
	IsProduction 	  bool
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
		//go func() {
			err := b.getPostsByGroupName(groupName)

			if err != nil {
				// TODO: error handling
				log.Println(err)
			}
		//}()
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	for index, post := range b.Resp[groupID].ResponseWall.Wall.Posts {
	//		commentID, err := b.getBestCommentOfPost(post.ID, groupID, index)
	//
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		b.Resp[groupID].ResponseWall.Wall.Posts[index].TopCommentID = commentID

			// b.Result = append(b.Result, ResultPost{
			// 	Body: b.Resp[groupID].ResponseWall.Wall.Posts[index].Text,
			// })
	//	}
	//}

	//log.Printf("%+v", b.Result)

	// for _, v := range b.Resp[-114383292].ResponseWall.Wall.Posts {
	// log.Printf("%+v", v)
	// log.Printf("%+v\n\n", b.Resp[-114383292].ResponseComments.CommentsList.Comments[v.ID])
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

func (b *Bot) getPostsByGroupName(groupName string) error {
	_, err := getGroupId(groupName, b.Config.AccessToken)

	if err != nil {
		return err
	}

	//getPostsByGroupIDURL := fmt.Sprintf(
	//	"%swall.get?owner_id=%d&count=10&v=5.52&access_token=%s",
	//	BaseAPIURL, groupId, b.Config.AccessToken)
	//
	//body, err := executeRequest(getPostsByGroupIDURL)
	//
	//if err != nil {
	//	return err
	//}
	//
	//log.Println(string(body))

	//b.Resp[groupName].ResponseWall = responseWall{}
///
//	json.Unmarshal(body, &b.Resp[groupName].ResponseWall)

	return nil
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
		"%sgroup.getById?group_id=%s&v=5.52&access_token=%s",
		BaseAPIURL, groupName, accessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return 0, err
	}

	log.Println(string(body))

	return 0, nil
}

//func (b *Bot) getBestCommentOfPost(postID int, groupID int, index int) (int, error) {
//	getCommentsByIDURL := fmt.Sprintf(
//		"%swall.getComments?owner_id=%d&post_id=%d&"+
//			"need_likes=1&count=100&v=5.52&access_token=%s",
//		BaseAPIURL, groupID, postID, b.AccessToken)
//
//	body, err := executeRequest(getCommentsByIDURL)
//
//	if err != nil {
//		return 0, err
//	}
//
//	b.Resp[groupID].ResponseComments = responseComments{}
//
//	// log.Println(string(body))
//
//	json.Unmarshal(body, &b.Resp[groupID].ResponseComments)
//
//	bestLikes := 0
//	id := 0
//
//	for _, comment := range b.Resp[groupID].ResponseComments.CommentsList.Comments {
//		if comment.Likes.Count > bestLikes {
//			bestLikes = comment.Likes.Count
//			id = comment.ID
//		}
//	}
//
//	return id, nil
//}
//
//func getAccessToken(clientID string, clientSecret string, b *Bot) error {
//	accessTokenURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s"+
//		"&client_secret=%s&v=5.80&grant_type=client_credentials",
//		clientID, clientSecret)
//
//	resp, err := http.Get(accessTokenURL)
//
//	if err != nil {
//		return err
//	}
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		return err
//	}
//
//	json.Unmarshal(body, &b)
//
//	return nil
//}

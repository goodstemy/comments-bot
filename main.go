package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// BaseAPIURL is base url for all requests
const BaseAPIURL = "https://api.vk.com/method/"

func main() {
	b := Bot{
		AccessToken: "7aec8e057aec8e057aec8e05607a887f8b77aec7aec8e0521dc3045ba1132f0568b38a4",
		Resp: map[int]*Response{
			-114383292: &Response{},
		},
	}

	// TODO: is that function needs?
	// err := getAccessToken("6615438", "WNOmK187BsnZedA1dEbX", &b)

	// if err != nil {
	// 	panic(err)
	// }

	for groupID := range b.Resp {
		err := b.getPostsByGroupID(groupID)

		if err != nil {
			panic(err)
		}

		for index, post := range b.Resp[groupID].ResponseWall.Wall.Posts {
			commentID, err := b.getBestCommentOfPost(post.ID, groupID, index)

			if err != nil {
				panic(err)
			}

			b.Resp[groupID].ResponseWall.Wall.Posts[index].TopCommentID = commentID

			// b.Result = append(b.Result, ResultPost{
			// 	Body:    b.Resp[groupID].ResponseWall.Wall.Posts[index].Text,
			// 	Comment: b.Resp[groupID].ResponseComments.CommentsList.Comments[].Text,
			// })
		}
	}

	log.Printf("%+v", b.Resp[-114383292])

	// for _, v := range b.Resp[-114383292].ResponseWall.Wall.Posts {
	// log.Printf("%+v", v)
	// log.Printf("%+v\n\n", b.Resp[-114383292].ResponseComments.CommentsList.Comments[v.ID])
	// }
}

func (b *Bot) getPostsByGroupID(groupID int) error {
	getPostsByGroupIDURL := fmt.Sprintf(
		"%swall.get?owner_id=%d&count=10&v=5.52&access_token=%s",
		BaseAPIURL, groupID, b.AccessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return err
	}

	b.Resp[groupID].ResponseWall = responseWall{}

	json.Unmarshal(body, &b.Resp[groupID].ResponseWall)

	return nil
}

func (b *Bot) getBestCommentOfPost(postID int, groupID int, index int) (int, error) {
	getCommentsByIDURL := fmt.Sprintf(
		"%swall.getComments?owner_id=%d&post_id=%d&"+
			"need_likes=1&count=100&v=5.52&access_token=%s",
		BaseAPIURL, groupID, postID, b.AccessToken)

	body, err := executeRequest(getCommentsByIDURL)

	if err != nil {
		return 0, err
	}

	b.Resp[groupID].ResponseComments = responseComments{}

	// log.Println(string(body))

	json.Unmarshal(body, &b.Resp[groupID].ResponseComments)

	bestLikes := 0
	id := 0

	for _, comment := range b.Resp[groupID].ResponseComments.CommentsList.Comments {
		if comment.Likes.Count > bestLikes {
			log.Println(comment.Likes.Count, comment.Text)
			bestLikes = comment.Likes.Count
			id = comment.ID
		}
	}

	return id, nil
}

func getAccessToken(clientID string, clientSecret string, b *Bot) error {
	accessTokenURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s"+
		"&client_secret=%s&v=5.80&grant_type=client_credentials",
		clientID, clientSecret)

	resp, err := http.Get(accessTokenURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	json.Unmarshal(body, &b)

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

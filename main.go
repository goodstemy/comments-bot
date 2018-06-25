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
	b := bot{
		AccessToken: "7aec8e057aec8e057aec8e05607a887f8b77aec7aec8e0521dc3045ba1132f0568b38a4",
		GroupID:     -114383292,
	}

	// err := getAccessToken("6615438", "WNOmK187BsnZedA1dEbX", &b)

	// if err != nil {
	// 	panic(err)
	// }

	err := b.getPostsByGroupID(b.GroupID)

	if err != nil {
		panic(err)
	}

	for _, post := range b.ResponseWall.Wall.Posts {
		err = b.getBestCommentOfPost(post.ID)

		if err != nil {
			panic(err)
		}
	}

	log.Println(b)
}

func (b *bot) getPostsByGroupID(groupID int) error {
	getPostsByGroupIDURL := fmt.Sprintf(
		"%swall.get?owner_id=%d&count=10&v=5.52&access_token=%s",
		BaseAPIURL, groupID, b.AccessToken)

	body, err := executeRequest(getPostsByGroupIDURL)

	if err != nil {
		return err
	}

	b.ResponseWall = response{}

	json.Unmarshal(body, &b.ResponseWall)

	return nil
}

func (b *bot) getBestCommentOfPost(ID int) error {
	getCommentsByIDURL := fmt.Sprintf(
		"%swall.getComments?owner_id=%d&post_id=%d&count=100&v=5.52&access_token=%s",
		BaseAPIURL, b.GroupID, ID, b.AccessToken)

	body, err := executeRequest(getCommentsByIDURL)

	if err != nil {
		return err
	}

	b.ResponseComments = response{}

	log.Println(string(body))

	json.Unmarshal(body, &b.ResponseComments)

	return nil
}

func getAccessToken(clientID string, clientSecret string, b *bot) error {
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

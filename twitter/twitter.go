package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterImage struct {
	Image_Type string
	W          int
	H          int
}

type MediaUploadResponse struct {
	Media_Id           int64
	Media_Id_String    string
	Media_Key          string
	Size               int64
	Expires_After_Secs int64
	Image              TwitterImage
}

func getClient() (twitter.Client, http.Client) {
	ConsumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	ConsumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	AccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	AccessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	token := oauth1.NewToken(AccessToken, AccessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter Client
	client := twitter.NewClient(httpClient)

	return *client, *httpClient
}

func Get(id int64) *twitter.Tweet {
	client, _ := getClient()

	tweet, _, _ := client.Statuses.Show(id, nil)
	return tweet

}

func Post(text string, imagePaths []string) *twitter.Tweet {

	if len(imagePaths) > 4 {
		log.Fatal("Too many images specified. Only 4 images may be attached to a Tweet")
	}

	// _, httpClient := getClient()
	client, httpClient := getClient()

	// Try to POST to upload/media
	base64, _ := ioutil.ReadFile("out-base64.txt")

	form := url.Values{}
	form.Add("media_data", string(base64))

	req, _ := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", strings.NewReader(form.Encode()))
	q := req.URL.Query()
	q.Add("media_category", "tweet_image")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var mur MediaUploadResponse
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &mur)

	fmt.Println(mur.Media_Id_String)

	tweet, resp, err := client.Statuses.Update(text, &twitter.StatusUpdateParams{
		MediaIds: []int64{mur.Media_Id},
	})
	log.Println(err)
	if resp != nil {
		log.Fatal(resp)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tweet)

	return tweet
}

// func uploadMedia(path string) {
// 	imageReader, imageReadErr := os.Open(path)
// 	if imageReadErr != nil {
// 		log.Fatal(imageReadErr)
// 	}

// 	client := http.Client{}
// 	req, err := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	// req.Form.Add("media_category", "tweet_image")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

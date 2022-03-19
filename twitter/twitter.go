package twitter

import (
	b64 "encoding/base64"
	"encoding/json"
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

/**
Upload a file to Twitter using the "POST media/upload" endpoint. The response includes a Media ID
that will be used later to associate the uploaded file to a new tweet. The Media ID is communicated
back to the caller using a channel

https://developer.twitter.com/en/docs/twitter-api/v1/media/upload-media/api-reference/post-media-upload
*/
func uploadMedia(imagePath string, client http.Client, c chan int64) int64 {
	log.Println(`Uploading image ` + imagePath)

	// Read the image file and convert it to base 64
	imgBuffer, imageReadErr := ioutil.ReadFile(imagePath)
	if imageReadErr != nil {
		log.Fatal(imageReadErr)
	}

	imageAsBase64 := b64.StdEncoding.EncodeToString(imgBuffer)

	// Create an HTTP form and set the media_data property
	form := url.Values{}
	form.Add("media_data", string(imageAsBase64))
	form.Add("media_category", "tweet_image")

	// Create a POST request that will upload the the accompanying image to the upload media endpoint
	req, createReqErr := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", strings.NewReader(form.Encode()))
	if createReqErr != nil {
		log.Fatal(createReqErr)
	}

	// Create a Query object that will set the media_category option. It's possible that

	// q := req.URL.Query()
	// q.Add("media_category", "tweet_image")
	// req.URL.RawQuery = q.Encode()

	// Add the Content-Type header to the request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Submit the POST request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Read the response from the upload request
	var mur MediaUploadResponse
	body, uploadResponseReadError := ioutil.ReadAll(res.Body)
	if uploadResponseReadError != nil {
		log.Fatal(uploadResponseReadError)
	}

	// Unmarshal the JSON body of the response into a MediaUploadResponse struct
	json.Unmarshal(body, &mur)

	// Send the MediaId over the channel. The listener will be able to use this
	// id to create a tweet with the image attached
	c <- mur.Media_Id

	// Return the MediaId as well, just in case something needs it
	return mur.Media_Id
}

func Post(text string, imagePaths []string) *twitter.Tweet {
	if len(imagePaths) > 4 {
		log.Fatal("Too many images specified. Only 4 images may be attached to a Tweet")
	}

	// Build a Twitter and a plain HTTP client configured with OAuth
	client, httpClient := getClient()
	mediaIdChannel := make(chan int64, len(imagePaths))

	// Upload all the media for this tweet
	for _, path := range imagePaths {
		go uploadMedia(path, httpClient, mediaIdChannel)
	}

	// Wait for the media IDs to come over the channel
	mediaIds := []int64{}
	for i := 0; i < len(imagePaths); i++ {
		mediaIds = append(mediaIds, <-mediaIdChannel)
	}

	// Post a tweet and attach all the uploaded media
	tweet, _, err := client.Statuses.Update(text, &twitter.StatusUpdateParams{
		MediaIds: mediaIds,
	})
	if err != nil {
		log.Fatal(err)
	}

	return tweet
}

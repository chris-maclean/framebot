package main

import (
	"flag"
	"log"

	"example.com/framebot/twitter"
)

type FrameBotState struct {
	File        string
	NextFrame   int
	TotalFrames int
}

func main() {
	log.Println("Hello, I am framebot!")
	// Locate the state file passed in as an argument
	stateFilePtr := flag.String("stateFile", "/home/chris/framebot-job.json", "Path to the framebot state file")

	flag.Parse()

	log.Println("state file located at: ", *stateFilePtr)

	// // Read and parse job file from local filesystem
	// bs, stateFileReadErr := ioutil.ReadFile(*stateFilePtr)
	// if stateFileReadErr != nil {
	// 	log.Fatal(stateFileReadErr)
	// }
	// state := FrameBotState{}

	// jsonErr := json.Unmarshal(bs, &state)
	// if jsonErr != nil {
	// 	log.Fatal(jsonErr)
	// }

	// log.Println("video file  : ", state.File)
	// log.Println("next frame  : ", state.NextFrame)
	// log.Println("total frames: ", state.TotalFrames)

	// t := twitter.Get(20)
	// log.Println(t.Text)

	twitter.Post("Hello world!", []string{})

	// Create image file from nth frame of video
	// ffmpeg.GetFrame(state.File, state.NextFrame)

	// Post tweet and attach image file
	// twitter.Post("Hello world!", nil)
	// Tweet info:
	// Current frame #
	// Total frame count
	// Name of film
	// Write job file to local filesystem
	// Exit

	// log.Println("Getting tweet with id=20")
	// tweet := twitter.Get(20)

	// log.Println(tweet.Text)

	// ffmpeg command to grab 35000th frame and save to a file
	// ffmpeg -v error -i 2001\ A\ Space\ Odyssey.mkv -vf "select=gte(n\,35000)" -vframes 1 out_img.jpg
	// or perhaps this one
	// ffmpeg -v error -i 2001\ A\ Space\ Odyssey.mkv -vf select='eq(n\,35000)' -vsync 0 out_img.jpg

}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chris-maclean/framebot/ffmpeg"
	"github.com/chris-maclean/framebot/twitter"
)

type FramebotState struct {
	File        string
	Title       string
	NextFrame   int
	TotalFrames int
}

func main() {
	// Locate the state file passed in as an argument
	// The default location of this file is   /opt/framebot/framebot-state.json
	stateFilePtr := flag.String("stateFile", "/opt/framebot/framebot-state.json", "Path to the framebot state file")

	flag.Parse()
	log.Println("state file located at: ", *stateFilePtr)

	// Extract the information from the state file, and apply some default values
	// so the state can be used going forward
	state := getState(stateFilePtr)

	// Log out information about the current run
	log.Println("Video file  : ", state.File)
	log.Println("Title       : ", state.Title)
	log.Println("Next frame  : ", state.NextFrame)
	log.Println("Total frames: ", state.TotalFrames)

	// Framebot can't pull the n-th frame from a movie with less than n frames. If the state
	// values say to do this, halt the program
	if state.NextFrame > state.TotalFrames {
		log.Fatalf("Cannot generate frame %d, only %d frames exist in %s. It's likely that Framebot has reached the end of the movie.",
			state.NextFrame,
			state.TotalFrames,
			state.Title)
		os.Exit(1)
	}

	// Build text of tweet
	text := buildText(state)
	log.Println(text)

	// Generate frame image and convert to base64
	framePath := "frame.jpeg"
	ffmpeg.GetFrame(state.File, state.NextFrame, framePath)

	// Post the tweet
	twitter.Post(text, []string{framePath})

	// Clean up the frame image
	deleteErr := os.Remove(framePath)
	if deleteErr != nil {
		log.Fatal(deleteErr)
	}

	// Update state to go to the next frame
	nextState := state
	nextState.NextFrame = state.NextFrame + 1

	// Write job file to the filesystem
	nextStateJson, marshalErr := json.Marshal(nextState)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	os.WriteFile(*stateFilePtr, nextStateJson, 0666)

	// Exit
	log.Printf("Tweeted frame %d of %s", state.NextFrame, state.Title)
	log.Println("Framebot complete!")

}

/**
Read the state file located at filePath. filePath will be either the value specified in the --stateFile option or the default location /opt/framebot/framebot-state.json
*/
func getState(filePath *string) FramebotState {
	// Read and parse job file from local filesystem
	bs, stateFileReadErr := os.ReadFile(*filePath)

	// If no file is found, halt the program
	if stateFileReadErr != nil {
		log.Fatal(stateFileReadErr)
	}

	// Create an empty state variable
	state := FramebotState{}

	// Unmarshal the bytes read from the file into a FramebotState struct
	jsonErr := json.Unmarshal(bs, &state)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// Assign default values if some properties are not present in the state file
	if state.File == "" {
		log.Printf("File location not set in state file. Setting path to /opt/framebot/movie")
		state.File = "/opt/framebot/movie"
	}

	if state.TotalFrames == 0 {
		log.Printf("Found 0 frames in state file, calculating total frames in %s", state.File)
		state.TotalFrames = ffmpeg.GetTotalFrames(state.File)
	}

	if state.NextFrame <= 0 {
		log.Printf("Next frame not set in state file. Setting NextFrame to 1")
		state.NextFrame = 1
	}

	return state
}

/**
Build the text of the tweet that will accompany the frame image
*/
func buildText(state FramebotState) string {
	// This line displays the current frame and the total number of frames
	// It will be the second line of text, if a Title is found
	line2 := fmt.Sprintf("Frame %d of %d", state.NextFrame, state.TotalFrames)

	// If no title is set in the state, don't add it to the text
	if state.Title == "" {
		return line2
	}

	// If there *is* a title in the state, make it the first line
	// of the tweet
	return fmt.Sprintf("%s\n%s", state.Title, line2)
}

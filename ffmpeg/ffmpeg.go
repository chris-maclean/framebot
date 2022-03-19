package ffmpeg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/**
Get the nth frame of a video file by shelling out to an ffmpeg command
https://stackoverflow.com/questions/20398539/extract-a-thumbnail-from-a-specific-video-frame
**/
func GetFrame(file string, n int, outFile string) {
	cmd := fmt.Sprintf("ffmpeg -v error -i '%s' -vf \"select=gte(n\\,%d)\" -vframes 1 %s", file, n, outFile)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated frame %d", n)
}

/**
https://ottverse.com/extract-frame-count-using-ffprobe-ffmpeg/
https://stackoverflow.com/questions/2017843/fetch-frame-count-with-ffmpeg
This call to ffprobe can be veeeeerrrryyyyy sloooowwwww, but it is only invoked if the state file doesn't contain a "totalFrames" property. If the user has a faster way to set this value, they should set it in the state file. That way, this function will never be called.
*/
func GetTotalFrames(file string) int {
	log.Printf("Calculating number of frames in %s", file)

	if _, fileExistsErr := os.Stat(file); errors.Is(fileExistsErr, os.ErrNotExist) {
		log.Fatalf("File %s does not exist", file)
	}

	cmd := fmt.Sprintf("ffprobe -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 '%s' | head -1", file)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	n, parseIntErr := strconv.Atoi(strings.TrimSpace(string(res)))
	if parseIntErr != nil {
		log.Fatal(parseIntErr)
	}

	log.Printf("Calculated %d total frames in %s", n, file)

	return n
}

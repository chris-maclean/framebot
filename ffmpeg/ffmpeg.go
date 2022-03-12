package ffmpeg

import (
	"log"
	"os/exec"
)

/**
* Get the nth frame of a video file by shelling out to an ffmpeg command
**/
func GetFrame(file string, n int) {
	out, err := exec.Command("bash", "-c", `which ffmpeg`).Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(out))
}

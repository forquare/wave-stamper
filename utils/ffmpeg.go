package utils

import (
	"log"
	"os/exec"
)

func ProcessVideo(image string, audio string, output string) bool {
	command := exec.Command("/usr/bin/ffmpeg", "-i", audio, "-i", image, "-filter_complex", "[0:a]aformat=channel_layouts=mono,showwaves=s=1080x600:mode=line:colors=orange[sw];[1][sw]overlay=(W-w)/2:(H-h)/2:format=auto,format=yuv420p[v]", "-map", "[v]", "-map", "0:a", "-vcodec", "libx265", "-crf", "28", output)
	out, err := command.Output()
	if err != nil {
		log.Println("ffmpeg failed")
		log.Println(out)
		log.Println(err)
		return false
	}

	return true
}

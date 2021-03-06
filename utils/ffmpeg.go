package utils

import (
	"log"
	"os/exec"
	"runtime"
)

// ProcessVideo returns a boolean on whether is was successful or not
func ProcessVideo(image string, audio string, output string) (int, string) {
	if !whereis("ffmpeg") {
		return 1, "Cannot find ffmpeg"
	}

	command := exec.Command("ffmpeg", "-i", audio, "-i", image, "-filter_complex", "[0:a]aformat=channel_layouts=mono,showwaves=s=1080x600:mode=line:colors=orange[sw];[1][sw]overlay=(W-w)/2:(H-h)/2:format=auto,format=yuv420p[v]", "-map", "[v]", "-map", "0:a", "-vcodec", "libx265", "-crf", "28", output)
	out, err := command.Output()
	if err != nil {
		log.Println("ffmpeg failed")
		log.Println(out)
		log.Println(err)
		return 2, "ffmpeg failed:\n" + string(out)
	}

	return 0, "done"
}

func whereis(name string) bool {
	var platform = runtime.GOOS
	var command = []string{"", "ffmpeg"}
	switch platform {
	case "windows":
		command[0] = "where"
		break
	default:
		command[0] = "which"
		break
	}

	cmd := exec.Command(command[0], command[1])
	err := cmd.Run()

	if err != nil {
		return false
	}
	return true
}

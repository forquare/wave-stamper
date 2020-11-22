package main

import (
	"fmt"
	"log"
	"os"

	// Do POSIX flags
	flag "github.com/spf13/pflag"

	// Local packages
	"github.com/forquare/wave-stamper/ui"
)

var (
	version = ""
	title   = "Wave Stamper"
	exeName = "wave-stamper"

	// Flags
	fVersion    bool
	fDebug      bool
	fAudioInput string
	fImageInput string
	fOutput     string
)

func initialiseFlags() {
	flag.BoolVarP(&fVersion, "version", "v", false, "Print the version number and quit.")
	flag.BoolVarP(&fDebug, "debug", "d", false, "Prints more details log information.")
	flag.StringVarP(&fAudioInput, "audio", "a", "", "mp3 file containing audio.")
	flag.StringVarP(&fImageInput, "image", "i", "", "Image file for background.")
	flag.StringVarP(&fOutput, "output", "o", "", "File to output movie to.")
}

func main() {
	initialiseFlags()

	flag.Parse()

	if flag.NArg()+flag.NFlag() == 0 {
		// No arguments/flags, launch GUI
		log.Println("Launching GUI")
		ui.GetUI(title, version)
	} else {
		if fVersion {
			fmt.Println(exeName, version)
			os.Exit(0)
		}

		if fDebug {
			log.Println("Command line mode")
		}
	}
}

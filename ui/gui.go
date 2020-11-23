package ui

import (
	"path/filepath"

	"github.com/andlabs/ui"
	// From github.com/andlabs/ui
	_ "github.com/andlabs/ui/winmanifest"

	"github.com/forquare/wave-stamper/utils"
)

var (
	mainwin     *ui.Window
	wintitle    string
	progversion string

	logopath  string
	audiopath string
	videopath string
)

func setupUI() {
	mainwin = ui.NewWindow(wintitle+" - "+progversion, 335, 200, true)

	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	leftvbox := ui.NewVerticalBox()
	leftvbox.SetPadded(true)
	hbox.Append(leftvbox, true)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	leftvbox.Append(grid, false)

	// Open Logo Controlls
	lblOpenLogo := ui.NewLabel("Logo File:")
	btnOpenLogo := ui.NewButton("...")
	txtLogo := ui.NewEntry()
	txtLogo.SetReadOnly(true)
	btnOpenLogo.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		ext := filepath.Ext(filename)
		if filename == "" {
			return
		}
		if ext != ".jpg" && ext != ".jpeg" {
			ui.MsgBoxError(mainwin,
				"Only JPEG images supported at this time",
				"Only files ending in '.jpg' or '.jpeg' are supported.")
			return
		}
		txtLogo.SetText(filename)
		logopath = filename
	})
	grid.Append(lblOpenLogo,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(txtLogo,
		0, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	grid.Append(btnOpenLogo,
		1, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	// Save Video Controlls
	lblSaveVideo := ui.NewLabel("Output File:")
	btnSaveVideo := ui.NewButton("...")
	txtVideo := ui.NewEntry()
	txtVideo.SetReadOnly(true)
	btnSaveVideo.OnClicked(func(*ui.Button) {
		filename := ui.SaveFile(mainwin)
		ext := filepath.Ext(filename)
		if filename == "" {
			return
		}
		if ext != ".mp4" && ext != ".mpeg4" {
			ui.MsgBoxError(mainwin,
				"Only MP4 and files are supported at this time",
				"Only files ending in '.mp4' or '.mpeg4' are supported.")
			return
		}
		txtVideo.SetText(filename)
		videopath = filename
	})
	grid.Append(lblSaveVideo,
		0, 5, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(txtVideo,
		0, 6, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	grid.Append(btnSaveVideo,
		1, 6, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	// Open Audio Controlls
	lblOpenAudio := ui.NewLabel("Audio File:")
	btnOpenAudio := ui.NewButton("...")
	txtAudio := ui.NewEntry()
	txtAudio.SetReadOnly(true)
	btnOpenAudio.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		ext := filepath.Ext(filename)
		base := filepath.Base(filename)[0 : len(filepath.Base(filename))-len(ext)]
		if filename == "" {
			return
		}
		if ext != ".mp3" && ext != ".mpeg3" && ext != ".wav" {
			ui.MsgBoxError(mainwin,
				"Only MP3 and wav files are supported at this time",
				"Only files ending in '.mp3' or '.mpeg3' or '.wav' are supported.")
			return
		}
		txtAudio.SetText(filename)
		audiopath = filename
		videopath = filepath.Clean(filepath.Dir(audiopath) + "/" + base + ".mp4")
		txtVideo.SetText(videopath)
	})
	grid.Append(lblOpenAudio,
		0, 2, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(txtAudio,
		0, 3, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	grid.Append(btnOpenAudio,
		1, 3, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	// Horizontal Separator
	grid.Append(ui.NewHorizontalSeparator(),
		0, 4, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	// SUBMIT!
	btnSubmit := ui.NewButton("Submit")
	btnSubmit.OnClicked(func(*ui.Button) {
		fileerror := false
		msg := ""

		if len(logopath) == 0 {
			msg += "You haven't specified a logo file.\n"
			fileerror = true
		}
		if len(audiopath) == 0 {
			msg += "You haven't specified an audio file.\n"
			fileerror = true
		}
		if len(videopath) == 0 {
			msg += "You haven't specified an output file.\n"
			fileerror = true
		}

		if fileerror {
			ui.MsgBoxError(mainwin,
				"Error", msg)
			return
		}
		utils.ProcessVideo(logopath, audiopath, videopath)
		ui.MsgBox(mainwin,
			"Done",
			"Completed video processing!")
	})
	grid.Append(btnSubmit,
		0, 7, 2, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	mainwin.SetChild(hbox)

	mainwin.Show()
}

// GetUI will take display and run the GUI
func GetUI(title string, version string) {
	wintitle = title
	progversion = version
	ui.Main(setupUI)
}

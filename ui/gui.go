package ui

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/forquare/wave-stamper/utils"
)

const (
	mainFormFilePath = "resources/ui.glade"
)

var (
	MainWin *MainWindow
)

var signals = map[string]interface{}{
	"on_logo_file_file_set":     on_logo_file_file_set,
	"on_mp3_file_file_set":      on_mp3_file_file_set,
	"on_find_save_file_clicked": on_find_save_file_clicked,
	"on_submit_clicked":         on_submit_clicked,
}

func ShowExistingMainWindow() *MainWindow {
	MainWin.Window.Show()
	MainWin.Window.Present()

	return MainWin
}

// This is a singleton
func GetMain(title string, version string) *MainWindow {
	if MainWin != nil && MainWin.Window.IsVisible() {
		return MainWin
	}

	MainWin = MainWindowNew(title, version)
	return MainWin
}

type MainWindow struct {
	Window *gtk.Window

	FcbLogoFile *gtk.FileChooserButton
	FcbMp3File  *gtk.FileChooserButton

	TxtFilePath *gtk.Entry

	ImgLogoImage *gtk.Image

	BtnFindFile *gtk.Button
	BtnSubmit   *gtk.Button

	PrgProgress *gtk.ProgressBar

	Title   string
	Version string

	LogoImagePath string
	Mp3FilePath   string
	OutputPath    string
}

func MainWindowNew(title string, version string) *MainWindow {
	win := new(MainWindow)
	win.Version = version
	win.Title = title

	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Fatal: %v", err)
		os.Exit(1)
	}

	err = builder.AddFromFile(mainFormFilePath)
	if err != nil {
		log.Fatal("Fatal: %v", err)
		os.Exit(1)
	}

	obj, err := builder.GetObject("main_window")
	if err != nil {
		log.Fatal("Fatal: %v", err)
		os.Exit(1)
	}

	builder.ConnectSignals(signals)

	var ok bool
	win.Window, ok = obj.(*gtk.Window)
	if !ok {
		log.Fatal("No main window found")
		os.Exit(1)
	}

	win.Window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.FcbLogoFile = utils.GetWidget(builder, "logo_file").(*gtk.FileChooserButton)
	win.ImgLogoImage = utils.GetWidget(builder, "logo_image").(*gtk.Image)
	win.FcbMp3File = utils.GetWidget(builder, "mp3_file").(*gtk.FileChooserButton)
	win.TxtFilePath = utils.GetWidget(builder, "save_file_text").(*gtk.Entry)
	win.BtnFindFile = utils.GetWidget(builder, "find_save_file").(*gtk.Button)
	win.BtnSubmit = utils.GetWidget(builder, "submit").(*gtk.Button)
	win.PrgProgress = utils.GetWidget(builder, "progress").(*gtk.ProgressBar)

	win.Window.SetTitle(win.Title)

	return win
}

// Handlers

func on_logo_file_file_set(fcb *gtk.FileChooserButton) {
	log.Println(fcb.GetFilename())

	MainWin.LogoImagePath = fcb.GetFilename()

	logoBuf, err := gdk.PixbufNewFromFileAtScale(MainWin.LogoImagePath, 250, 250, true)
	if err != nil {
		log.Fatal("Fatal: %v", err)
		os.Exit(1)
	}

	MainWin.ImgLogoImage.SetFromPixbuf(logoBuf)
}

func on_mp3_file_file_set(fcb *gtk.FileChooserButton) {
	log.Println(fcb.GetFilename())

	MainWin.Mp3FilePath = fcb.GetFilename()
}

func on_find_save_file_clicked(btn *gtk.Button) {
	log.Println("Go get file!")
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Output File",
		MainWin.Window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Cancel",
		gtk.RESPONSE_CANCEL,
		"Save",
		gtk.RESPONSE_ACCEPT)

	// Only MP4s are currently supported as output
	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.mp4")
	dialog.SetFilter(filter)

	// If we already have a string, let's dump the user in the same directory
	if len(MainWin.OutputPath) > 0 {
		dialog.SetCurrentFolder(filepath.Dir(MainWin.OutputPath))
		dialog.SetCurrentName(filepath.Base(MainWin.OutputPath))
	}

	res := dialog.Run()
	if res != gtk.RESPONSE_ACCEPT {
		dialog.Destroy()
		return
	}

	MainWin.OutputPath = dialog.GetFilename()
	MainWin.TxtFilePath.SetText(MainWin.OutputPath)
	log.Println(MainWin.OutputPath)

	dialog.Destroy()

	// Check that the output file is going to be an MP4
	var mp4regex = regexp.MustCompile(`.*\.mp4$`)
	if !mp4regex.MatchString(MainWin.OutputPath) {
		fileError := gtk.MessageDialogNew(
			MainWin.Window,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_ERROR,
			gtk.BUTTONS_OK,
			"Output file must be an mp4.")
		fileError.Run()
		fileError.Destroy()
		on_find_save_file_clicked(btn)
		return
	}
}

func on_submit_clicked(btn *gtk.Button) {
	log.Println("Push the button!")
	// Make sure any edits the user has made are picked up
	MainWin.OutputPath, _ = MainWin.TxtFilePath.GetText()
	log.Println(MainWin.OutputPath)
	utils.ProcessVideo(MainWin.LogoImagePath, MainWin.Mp3FilePath, MainWin.OutputPath)
}

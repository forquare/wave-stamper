package utils

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func GetWidget(builder *gtk.Builder, id string) glib.IObject {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Fatal("Fatal: %v", err)
		os.Exit(1)
	}
	return obj
}

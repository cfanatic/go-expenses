package main

import (
	"log"
	"os"

	"github.com/cfanatic/go-expensegui/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	DIR = "/Users/cfanatic/Coding/Go/src/github.com/cfanatic/go-expense/res/expenses/"
)

func main() {
	var app = widgets.NewQApplication(len(os.Args), os.Args)
	var win = gui.NewGui(nil, 0)

	if DIR[len(DIR)-1:] != "/" {
		log.Println("Terminate constant 'DIR' with a forward slash")
		return
	}

	win.InitWith(app)
	win.Show()
	app.Exec()
}

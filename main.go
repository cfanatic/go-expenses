package main

import (
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

	win.InitWith(app)
	win.Show()
	app.Exec()
}

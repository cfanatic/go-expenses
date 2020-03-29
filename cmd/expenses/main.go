package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cfanatic/go-expenses/account"
	"github.com/cfanatic/go-expenses/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	MODE = "GUI"
	DIR  = "/Users/cfanatic/Coding/Go/src/github.com/cfanatic/go-expenses/misc/transactions/"
)

func main() {
	switch MODE {
	case "CLI":
		var acc account.IAccount
		var exps []*account.Expense
		if DIR[len(DIR)-1:] != "/" {
			log.Println("Terminate constant 'DIR' with a forward slash")
			return
		}
		if dir, _ := ioutil.ReadDir(DIR); len(dir) > 0 {
			for _, file := range dir {
				str := strings.Split(file.Name(), ".")
				if len(str) == 1 || strings.Contains(str[0], "~$") || str[1] != "xlsx" {
					continue
				}
				path := DIR + file.Name()
				acc = &account.Expense{Path: path}
				acc.Init()
				acc.Run()
				acc.Print()
				exps = append(exps, acc.(*account.Expense))
			}
			acc = &account.Expenses{Exp: exps}
			acc.Init()
			acc.Run()
			acc.Print()
			acc.Plot()
		} else {
			log.Println("Excel sheets not found")
		}
	case "GUI":
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
}

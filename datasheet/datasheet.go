package datasheet

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	DATE1 = iota
	DATE2
	PAYEE
	TYPE
	DESC
	AMOUNT
)

type Datasheet struct {
	path string
	tab  string
	file *excelize.File
	Err  error
}

type Content struct {
	Date   string
	Payee  string
	Desc   string
	Amount float32
}

func New(path, tab string) *Datasheet {
	ds := Datasheet{path: path, tab: tab}
	if f, err := excelize.OpenFile(path); err == nil {
		ds.file, ds.Err = f, nil
		return &ds
	} else {
		ds.path, ds.Err = "", err
		return &ds
	}
}

func (ds *Datasheet) Cell(col string, row int) (string, error) {
	element := col + strconv.Itoa(row)
	cell, err := ds.file.GetCellValue(ds.tab, element)
	return cell, err
}

func (ds *Datasheet) Content() ([]Content, error) {
	cont := make([]Content, 0)
	if rows, err := ds.file.GetRows(ds.tab); err == nil {
		for _, row := range rows {
			date := row[DATE1]
			if len(date) == 8 {
				// Format given in MM-DD-YY already
			} else if len(date) == 10 {
				// Convert format from YYYY-MM-DD to MM-DD-YY
				tmp := strings.Split(date, "-")
				date = fmt.Sprintf("%s-%s-%s", tmp[1], tmp[2], tmp[0][2:])
			} else {
				panic("Invalid date format provided by Datasheet")
			}
			amount, _ := strconv.ParseFloat(row[AMOUNT], 32)
			data := Content{date, row[PAYEE], row[DESC], float32(amount)}
			cont = append(cont, data)
		}
		return cont, nil
	} else {
		return cont, err
	}
}

func (ds *Datasheet) Update(old, new Content) error {
	var date, payee, desc, amount []string
	var ch chan []string

	// Horrible workaround due to floating point rounding errors in Excelize
	// The problem is that -40.8 is internally stored as -40.799999999999997
	// Solution: -40.8 is converted into 40 and used as a regex search parameter
	tmp := fmt.Sprintf("\\b%.0f\\b", math.Abs(math.Trunc((float64(old.Amount)))))

	// Convert the date format if YYYY-MM-DD is used in the Excel sheet
	cell, _ := ds.Cell("A", 1)
	if len(old.Date) == 8 && len(cell) == 10 {
		tmp := strings.Split(old.Date, "-")
		old.Date = fmt.Sprintf("%s-%s-%s", "20"+tmp[2], tmp[0], tmp[1])
	}

	ch = make(chan []string)
	cut := func(list []string, ch chan []string) {
		slice := []string{}
	LOOP:
		for _, item := range list {
			for _, val := range slice {
				if val == item[1:len(item)] {
					continue LOOP
				}
			}
			slice = append(slice, item[1:len(item)])
		}
		ch <- slice
	}

	date, _ = ds.file.SearchSheet(ds.tab, old.Date)
	payee, _ = ds.file.SearchSheet(ds.tab, old.Payee)
	desc, _ = ds.file.SearchSheet(ds.tab, old.Desc)
	amount, _ = ds.file.SearchSheet(ds.tab, tmp, true)

	go cut(date, ch)
	go cut(payee, ch)
	go cut(desc, ch)
	go cut(amount, ch)

	date, payee, desc, amount = <-ch, <-ch, <-ch, <-ch

	if len(date) > 0 && len(payee) > 0 && len(desc) > 0 && len(amount) > 0 {
		indices := []string{}
		indices = append(indices, date...)
		indices = append(indices, payee...)
		indices = append(indices, desc...)
		indices = append(indices, amount...)
		sort.Strings(indices)
		cnt, row := 0, ""
	LOOP:
		for i := 0; i < len(indices)-1; i++ {
			if indices[i] == indices[i+1] {
				cnt++
				switch {
				case cnt == 1:
					row = indices[i]
				case cnt == 3:
					break LOOP
				}
			} else {
				cnt = 0
				row = ""
			}
		}
		cellp := string('A'+PAYEE) + row
		cella := string('A'+AMOUNT) + row
		ds.file.SetCellStr(ds.tab, cellp, new.Payee)
		ds.file.SetCellFloat(ds.tab, cella, float64(new.Amount), 2, 32)
		ds.file.Save()
		return nil
	} else {
		return errors.New("Error: Unable to find datasheet row index")
	}
}

func (ds *Datasheet) Print(content interface{}) {
	log.Printf("%+v\n", content)
}

package main

import (
	//	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	//	"strconv"
	"strings"
)

////------------ Объявление типов и глобальных переменных

var (
	hd   string
	user string
)

type page struct {
	Title string
	Msg   string
}

//------------ END Объявление типов и глобальных переменных

// сохранить файл
func Savestrtofile(namef string, str string) int {
	file, err := os.Create(namef)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	title := r.URL.Path[len("/"):]

	if title != "exec/" {
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, &page{Title: "Создание файла для рассыли почты", Msg: "Задание триггера (условия) на срабатывание бота цен"})
	} else {
		mstr := r.FormValue("multistroka")
		mm := strings.Split(mstr, "\r\n")
		fmt.Println("Длина: ", len(mm))

		tm := make([]string, 0)

		sres := ""

		for i := 0; i < len(mm); i++ {
			if mm[i] != "" {
				s := mm[i]
				tm = append(tm, s)
				sres += s + "\n"
			}
		}

		Savestrtofile("pochta.cfg", sres)

		t1, _ := template.ParseFiles("template-result.html")
		t1.Execute(w, &page{Title: "Введенные данные: \n ", Msg: sres})

	}
}

func main() {

	http.HandleFunc("/", index)

	http.ListenAndServe(":7777", nil)
}

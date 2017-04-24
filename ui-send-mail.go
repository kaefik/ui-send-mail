package main

import (
	//	"flag"
	"bytes"
	"fmt"
	//	"html/template"
	"net/http"
	"os"
	"os/exec"
	//	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"
)

////------------ Объявление типов и глобальных переменных

var (
	hd   string
	user string
)

var (
	tekuser     string // текущий пользователь который задает условия на срабатывания
	pathcfg     string // адрес где находятся папки пользователей, если пустая строка, то текущая папка
	pathcfguser string
)

type page struct {
	Title  string
	Msg    string
	Msg2   string
	TekUsr string
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

func indexHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {

	rr.HTML(200, "index", &page{Title: "Работа с конфиг файлом лога звонков", Msg: "Начальная страница", TekUsr: "Текущий пользователь: " + string(user)})
}

// обработка редактирования списка получателей
func EditRecipientHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	//	nstr, _ := strconv.Atoi(params["nstr"])
	//	tt := dataConfigFile[nstr]
	rr.HTML(200, "editrecipient", "")
}

// обработка сохранения списка получателей
func SaveRecipientsHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {

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
	sres += "n.zaharova@kazan.2gis.ru" + "\n"
	sres += "i.saifutdinov@kazan.2gis.ru" + "\n"

	Savestrtofile("pochta.cfg", sres)

	rr.HTML(200, "saverecipients", "")
}

// обработка редактирования списка получателей
func EditShablonHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	//	nstr, _ := strconv.Atoi(params["nstr"])
	//	tt := dataConfigFile[nstr]
	rr.HTML(200, "editshablon", "")
}

// обработка сохранения шаблона письма
func SaveShablonHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
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
	Savestrtofile("shablonmail", sres)
	rr.HTML(200, "saveshablon", "")
}

// отправка писем получателям
func SendMailHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	cmd := exec.Command("/bin/sh", "client-sendmail.sh")
	//	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())

	rr.HTML(200, "sendmail", "")
}

func authFunc(username, password string) bool {
	return (auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "qwe123!!"))
}

func main() {

	fmt.Println("Starting...")
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates", // Specify what path to load the templates from.
		Layout:     "layout",    // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}}))

	m.Use(auth.BasicFunc(authFunc))
	m.Get("/editshablon", EditShablonHandler)
	m.Get("/sendmail", SendMailHandler)
	m.Post("/saveshablon", SaveShablonHandler)
	m.Get("/editrecipient", EditRecipientHandler)
	m.Post("/saverecipients", SaveRecipientsHandler)
	m.Get("/", indexHandler)
	m.RunOnAddr(":7777")
}

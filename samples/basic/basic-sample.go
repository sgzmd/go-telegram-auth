package main

import (
	"flag"
	"go-telegram-auth/tgauth"
	"html/template"
	"log"
	"net/http"
)

const (
	TelegramCookie = "tg_auth"

	CheckAuthPage = "/check-auth"
	AuthPage      = "/auth"
	DefaultPort   = "8080"

	// Domain TODO: change to the domain you configured (see the docs)
	Domain = "tgauth.com"

	// BotName TODO: CHANGE THIS TO YOUR OWN
	BotName = "sgzmd_tgauth_bot"

	Html = `<!DOCTYPE html>
<html><head><title>Go Web Server</title></head>
<body><h1>Go Web Server</h1>
<h1>Hello, anonymous!</h1>
<script async src="https://telegram.org/js/telegram-widget.js?19" data-telegram-login="{{.BotName}}" data-size="large" data-auth-url="http://{{.Domain}}/{{.CheckAuthUrl}}" data-request-access="write"></script>
</body></html>`
)

var TgAuthKey string
var Auth tgauth.TelegramAuth

func main() {
	tgapi := flag.String("telegram_api_key", "", "Telegram API key")
	flag.Parse()

	if *tgapi == "" {
		panic("Telegram API key is required")
	}

	TgAuthKey = *tgapi
	Auth = tgauth.NewTelegramAuth(TgAuthKey, AuthPage, CheckAuthPage)

	http.HandleFunc(CheckAuthPage, HandleAuth)
	http.HandleFunc(AuthPage, HandleLoginPage)
	http.HandleFunc("/", HandleIndexPage)

	e := http.ListenAndServe("tgauth.com:"+DefaultPort, nil)
	if e != nil {
		panic(e)
	}
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := make(map[string][]string)
	for k, v := range r.Form {
		params[k] = v
	}

	ok, err := Auth.CheckAuth(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Invalid auth", http.StatusUnauthorized)
		return
	}

	Auth.SetCookie(w, params)

	// redirect back to the main page
	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleLoginPage(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.New("index").Parse(Html))
	tmpl.Execute(writer, struct {
		BotName      string
		Domain       string
		CheckAuthUrl string
	}{BotName, Domain, CheckAuthPage})
}

func HandleIndexPage(writer http.ResponseWriter, request *http.Request) {
	params, err := Auth.GetParamsFromCookie(request)
	if err != nil {
		log.Printf("Unable to get params from cookie: %+v", err)
		http.Redirect(writer, request, "/auth", http.StatusFound)
		return
	}

	ok, err := Auth.CheckAuth(params)
	if err != nil {
		log.Printf("Unable to check auth: %+v", err)
		http.Redirect(writer, request, "/auth", http.StatusFound)
		return
	} else if !ok {
		log.Printf("Auth is not ok")
		http.Redirect(writer, request, "/auth", http.StatusFound)
		return
	}

	writer.Write([]byte("<html><body><h1>Welcome, " + params["first_name"][0] + "</h1></body></html>"))
}

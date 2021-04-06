package apiserver

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iavealokin/cashflow/app/store"
	"github.com/iavealokin/cashflow/app/model"
	"github.com/sirupsen/logrus"
)




type webserver struct {
	router 		 *mux.Router
	logger 		 *logrus.Logger
	store 		 store.Store
}


func newWebServer(store store.Store) *webserver{
ws :=&webserver{
	router: 	  mux.NewRouter(),
	logger: 	  logrus.New(),
	store: 		  store,
}
ws.router.StrictSlash(true)
staticDir := "/public/"
ws.router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
ws.configureWebRouter()
return ws
}

func (ws *webserver) ServeHTTP(w http.ResponseWriter, r *http.Request){
	ws.router.ServeHTTP(w,r)
}

func (ws *webserver) configureWebRouter(){
	ws.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	staticDir := "/public/"
	http.FileServer(http.Dir(staticDir))
	ws.router.HandleFunc("/",ws.handleStartPage)
	ws.router.HandleFunc("/home", ws.handleUserLogin)
//	ws.router.HandleFunc("/transactions",ws.handleUserTransactions)
}


func (ws*webserver) handleStartPage(w http.ResponseWriter, r *http.Request){
	//указываем путь к нужному файлу
	path := filepath.Join("public", "html", "login.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//выводим шаблон клиенту в браузер
	err = tmpl.ExecuteTemplate(w,"error","")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (ws*webserver) handleUserLogin(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		login:= r.FormValue("login")
		password :=r.FormValue("password")
		user,err:=ws.store.User().UserLogin(login,password)
		if err!=nil{
			path := filepath.Join("public", "html", "login.html")
			//создаем html-шаблон
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		
			//выводим шаблон клиенту в браузер
			err = tmpl.ExecuteTemplate(w,"error","Incorrect login or password!")
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			
		}else{
			operations,err := ws.store.User().GetOperations(user.ID); 
			usrData, err :=ws.store.User().GetUserData(user.ID);
	if err != nil{
		ws.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	//указываем путь к нужному файлу
	main := filepath.Join("public","html","template.html")
	common := filepath.Join("public","html","home.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(main,common)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
type TemplateCustom struct{
	Username 	string
	Uname 		string
	Income 		string
	Outcome 	string
	Difference  string
	Flag 		int 
}
 //tc:= new(TemplateCustom)
 tc := new(model.UserData)
 tc.Username	= user.Surname+" " +user.Username
 tc.Uname 		= string(user.Surname[0])+string(user.Username[0])
 tc.Income 		= usrData.Income
 tc.Outcome 	= usrData.Outcome
 tc.Difference  = usrData.Difference
 tc.Flag        = usrData.Flag
 tc.Actives 	= usrData.Actives
 tc.Passives	= usrData.Passives
 tc.Categories	= usrData.Categories
 tc.Operations  = usrData.Operations
	//выводим шаблон клиенту в браузер
	err = tmpl.ExecuteTemplate(w,"operations",struct{Operations, User interface{}}{operations,tc})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
		}
			
}


func (ws *webserver) error(w http.ResponseWriter, r *http.Request, code int, err error){
	ws.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (ws *webserver) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}){
	w.WriteHeader(code)
	if data != nil{
		json.NewEncoder(w).Encode(data)
	}
}

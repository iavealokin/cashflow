package model
//User struct ...
type User struct {
	ID       int    `json:"userid"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Surname  string `json:"surname"`
	Birthday string `json:"birthday"`
	Password string `json:"password"`
}
type Operation struct { 
	ID        string    `json:"operationid"`
	Amount    string `json:"amount"`
	Direction string `json:"direction"`
	Comment   string `json:"comment"`
	Date 	  string `json:"date"`

}

type UserActive struct{
	ID              string
 	Name    		string
 	Cost            string
 	Amount          string
 	Result          string
 	Percent         string
	UserID	        string
	Sum				string
}
type UserPassive struct{
	ID              string
 	Name    		string
 	Cost            string
 	Amount          string
 	Result          string
 	Percent         string
	UserID	        string
	Sum				string
}

type UserData struct {
    Username 	string
	Uname 		string
	Income 		string
	Outcome 	string
	Difference  string
	Flag 		int
    Actives		[] UserActive
    Passives	[] UserPassive
}
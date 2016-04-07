package route

import (
	"corate/util"
	"corate/model"
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
)

func DashboardHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session,_:=util.GlobalSessions.SessionStart(w,r)
	defer session.SessionRelease(w)

	files:=[]string{"base","dashboard"}
	user:=session.Get("user").(map[string]interface{})
	quotes:=model.GetQuotesByUser(user["id"].(string))
	for _,q:=range quotes{
		
		// Calculate time since highlighted
		q["created_at"]=util.TimeShow(q["created_at"].(float64))

		// Tab according to keywords
		keywords:=q["keywords"].([]interface{})
		tab:=map[string]bool{"startup":false,"personal":false,"work":false}
		for _,k:=range keywords{
			text:=strings.ToUpper(k.(map[string]interface{})["text"].(string))
			if strings.Index(text,"STARTUP")>=0 ||
				strings.Index(text,"TECH")>=0{
				tab["startup"]=true
			}
			if strings.Index(text,"PERSONAL")>=0 ||
				strings.Index(text,"PRODUCTIVITY")>=0 ||
				strings.Index(text,"LIFESTYLE")>=0{
				tab["personal"]=true
			}
		}
		q["tab"]=tab
	}
	data:=make(map[string]interface{},2)
	data["user"]=user
	data["list"]=quotes
	util.RenderTemplate(w,data,files...)
}

func SendUserHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session,_:=util.GlobalSessions.SessionStart(w,r)
	defer session.SessionRelease(w)

	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Content-Type","application/json")

	var resp []byte
	var err error
	if session.Get("user")!=nil{
		user:=session.Get("user").(map[string]interface{})
		user["authenticated"]=1
		resp,err=json.Marshal(user)
		
	} else{
		resp,err=json.Marshal(map[string]interface{}{"authenticated":0})
	}
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func LogoutHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)

	util.GlobalSessions.SessionDestroy(w,r)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
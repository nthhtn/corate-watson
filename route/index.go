package route

import (
	"corate/util"
	"net/http"
	"fmt"
)

func IndexHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session,_:=util.GlobalSessions.SessionStart(w,r)
	defer session.SessionRelease(w)

	if session.Get("profile")!=nil{
		http.Redirect(w,r,"/dashboard",http.StatusMovedPermanently)
	} else{
		files:=[]string{"base","index"}
		util.RenderTemplate(w,nil,files...)
	}
}

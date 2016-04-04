package route

import (
	"corate/util"
	"net/http"
	"fmt"
)

func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)
	
	if session.Get("profile") == nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		next(w, r)
	}
}

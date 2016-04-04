package route

import (
	"corate/util"
	"corate/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

var conf *oauth2.Config

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method," ",r.URL)
	conf=&oauth2.Config{
		ClientID:"700740834863-m4om9r91htn19htq2b6a05fu6vu4j7i5.apps.googleusercontent.com",
		ClientSecret:"93QSpD0qbgbwGQsZe934s-rB",
		RedirectURL:"http://localhost:3000/auth/google",
		Scopes:[]string{"profile","email"},
		Endpoint:google.Endpoint,
	}
	http.Redirect(w, r, conf.AuthCodeURL("state"), http.StatusMovedPermanently)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method," ",r.URL)

	// Exchange for profile
	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user session
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)

	session.Set("id_token", token.Extra("id_token"))
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)

	users:=model.GetUsersByField("email",profile["email"])
	if (len(users)>0){
		u:=users[0]
		u["name"]=profile["name"].(string)
		u["avatar"]=profile["picture"].(string)
		model.UpdateUser(u["id"].(string),u)
		session.Set("user",model.GetUsersByField("id",u["id"])[0])
	} else{
		u:=map[string]interface{}{
			"name":profile["name"].(string),
			"email":profile["email"].(string),
			"token":session.Get("access_token").(string),
			"avatar":profile["picture"].(string),
			"type":"google",
		}
		idU:=model.InsertUser(u).GeneratedKeys[0]
		session.Set("user",model.GetUsersByField("id",idU)[0])
	}

	http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
}
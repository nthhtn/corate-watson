package main

import (
	"corate/config"
	"corate/route"
	"corate/util"
	"github.com/codegangsta/negroni"
	"net/http"
)

func main(){
	
	// Setup database
	config.Setupdb()

	// Initialize utilities
	util.InitSessionManager()

	// Serve static files
	http.Handle("/public/",http.StripPrefix("/public/",http.FileServer(http.Dir("public"))))

	// Setup routers and middlewares

	// Index route
	http.HandleFunc("/",route.IndexHandler)

	// User routes
	http.HandleFunc("/login",route.GoogleLoginHandler)
	http.HandleFunc("/auth/google",route.GoogleCallbackHandler)
	http.HandleFunc("/api/oauth",route.SendUserHandler)
	http.Handle("/dashboard",negroni.New(
		negroni.HandlerFunc(route.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(route.DashboardHandler)),
	))

	// Quote routes
	http.Handle("/api/on",negroni.New(
		negroni.HandlerFunc(route.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(route.SendQuotesHandler)),
	))
	http.Handle("/api/create",negroni.New(
		negroni.HandlerFunc(route.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(route.SaveQuoteHandler)),
	))
	http.Handle("/tag",negroni.New(
		negroni.HandlerFunc(route.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(route.TagHandler)),
	))
	http.Handle("/search",negroni.New(
		negroni.HandlerFunc(route.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(route.TagHandler)),
	))

	// Start listening
	http.ListenAndServe(":3000",nil)
}
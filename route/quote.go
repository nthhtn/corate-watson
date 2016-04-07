package route

import (
	"corate/util"
	"corate/model"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"strings"
	"time"
	"fmt"
)

func SendQuotesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)
	
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Content-Type","application/json")

	// Decode request body
	raw,err:=ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	data:=struct{
		URL string `json:"url"`
	}{}
	err=json.NewDecoder(bytes.NewReader(raw)).Decode(&data)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	// Send highlighted texts to client
	idUser:=session.Get("user").(map[string]interface{})["id"].(string)
	quotes:=model.GetQuotesByUrl(data.URL,idUser)
	var resp []byte
	if (len(quotes)>0){
		text:=make([]map[string]interface{},0)
		for _,quote:=range quotes{
			text=append(text,map[string]interface{}{
				"id":quote["id"],
				"text":quote["htmltext"],
				"path":quote["path"],
			})
		}
		resp,err=json.Marshal(map[string]interface{}{
			"found":1,
			"text":text,
		});
	} else{
		resp,err=json.Marshal(map[string]interface{}{"found":0})
	}
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func SaveQuoteHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)
	
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Content-Type","application/json")

	var resp []byte
	var err error
	url:=r.FormValue("url")
	quote:=map[string]interface{}{
		"text":r.FormValue("text"),
		"webtitle":r.FormValue("title"),
		"weburl":url,
		"path":r.FormValue("nodePath"),
		"htmltext":r.FormValue("htmltext"),
		"created_at":time.Now().UTC().UnixNano()/1e6,
		"idUser":r.FormValue("id"),
	}

	// Insert new quote
	res:=model.InsertQuote(quote)
	if res.Inserted==1{
		fmt.Println("Insert quote successfully")
		idQuote:=res.GeneratedKeys[0]
		resp,err=json.Marshal(map[string]interface{}{
			"error":0,
			"id":idQuote,
		})
		if err!=nil{
			fmt.Println(err)
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		articles:=model.GetArticlesByField("url",url)

		if (len(articles)==0){

			// Insert new article
			article:=map[string]interface{}{
				"url":url,
				"keywords":[]map[string]interface{}{},
				"taxonomy":[]map[string]interface{}{},
				"concepts":[]map[string]interface{}{},
				"author":"",
			}
			res=model.InsertArticle(article)
			var idArticle string
			if res.Inserted==1{
				fmt.Println("Insert article successfully")
				idArticle=res.GeneratedKeys[0]

				// Update quote with created Article ID
				res=model.UpdateQuote(idQuote,map[string]interface{}{
					"idArticle":idArticle,
				})

				// Analyze AlchemyAPI
				util.InitAnalyzer()

				// Keywords
				keywords:=util.Keywords(url)
				res=model.UpdateArticle(idArticle,map[string]interface{}{
					"keywords":keywords,
				})
				if res.Errors==0{
					fmt.Println("Update keywords")
				}

				// Taxonomy
				taxonomy:=util.Taxonomy(url)
				res=model.UpdateArticle(idArticle,map[string]interface{}{
					"taxonomy":taxonomy,
				})
				if res.Errors==0{
					fmt.Println("Update taxonomy")
				}

				// Concepts
				concepts:=util.Concepts(url)
				res=model.UpdateArticle(idArticle,map[string]interface{}{
					"concepts":concepts,
				})
				if res.Errors==0{
					fmt.Println("Update concepts")
				}

				// Author
				authors:=util.Concepts(url)
				if (len(authors)>0){
					res=model.UpdateArticle(idArticle,map[string]interface{}{
					"author":authors[0],
					})
					if res.Errors==0{
						fmt.Println("Update author")
					}
				}
			}
		} else{

			// Update quote with existing Article ID
			res=model.UpdateQuote(idQuote,map[string]interface{}{
				"idArticle":articles[0]["id"],
			})
		}
	} else{
		resp,err=json.Marshal(map[string]interface{}{
			"error":1,
			"message":"Cannot insert quote",
		})
		if err!=nil{
			fmt.Println(err)
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}
}

func DeleteQuoteHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)
	
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Content-Type","application/json")

	// Decode request body
	raw,err:=ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	data:=struct{
		IdQuote string `json:"idQ"`
	}{}
	err=json.NewDecoder(bytes.NewReader(raw)).Decode(&data)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	// Delete quote
	res:=model.DeleteQuote(data.IdQuote)
	resp,err:=json.Marshal(map[string]interface{}{"deleted":res.Deleted})
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func TagHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method," ",r.URL)
	session, _ := util.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)

	files:=[]string{"base","tag"}
	user:=session.Get("user").(map[string]interface{})
	var tag string
	if strings.Index(r.URL.String(),"/tag")==0{
		tag=r.URL.Query().Get("tag")
	} else{
		tag=r.URL.Query().Get("k")
	}
	quotes:=model.GetQuotesByTag(tag,user["id"].(string))
	for _,q:=range quotes{
		q["created_at"]=util.TimeShow(q["created_at"].(float64))
	}
	data:=make(map[string]interface{},3)
	data["user"]=user
	data["list"]=quotes
	data["tag"]=tag
	util.RenderTemplate(w,data,files...)
}
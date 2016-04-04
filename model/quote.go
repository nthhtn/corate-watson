package model

import (
	"corate/config"
	r "github.com/dancannon/gorethink"
	"fmt"
)

func InsertQuote(quote interface{}) r.WriteResponse{
	result,err:=r.Table("quote").Insert(quote).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func UpdateQuote(id string,data interface{}) r.WriteResponse{
	result,err:=r.Table("quote").Get(id).Update(data).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func GetQuotesByUser(idUser string) []map[string]interface{}{
	result,err:=r.Table("quote").Filter(r.Row.Field("idUser").Eq(idUser)).
		OrderBy(r.Desc("created_at")).
		Merge(func(quote r.Term) interface{}{
			article:=r.Table("article").Get(quote.Field("idArticle"))
			return map[string]interface{}{
				"keywords":article.Field("keywords"),
			}
		}).Run(config.Connection())
	defer result.Close()
	
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	if result.IsNil(){
		return []map[string]interface{}{}
	}
	var all []interface{}
	err=result.All(&all)
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	quotes:=make([]map[string]interface{},len(all))
	for i,cur:=range all{
		quotes[i]=cur.(map[string]interface{})
	}
	return quotes
}

func GetQuotesByUrl(url string,idUser string) []map[string]interface{}{
	result,err:=r.Table("quote").Filter(map[string]interface{}{
		"weburl":url,
		"idUser":idUser,
	}).Run(config.Connection())
	defer result.Close()

	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	if result.IsNil(){
		return []map[string]interface{}{}
	}
	var all []interface{}
	err=result.All(&all)
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	quotes:=make([]map[string]interface{},len(all))
	for i,cur:=range all{
		quotes[i]=cur.(map[string]interface{})
	}
	return quotes
}

func GetQuotesByTag(tag string,idUser string) []map[string]interface{}{
	articles:=GetArticlesByTag(tag)
	if (len(articles)>0){
		ids:=make([]interface{},len(articles))
		for i,article:=range articles{
			ids[i]=article["id"]
		}
		result,err:=r.Table("quote").Filter(func (quote r.Term) r.Term{
				return r.Expr(ids).Contains(quote.Field("idArticle"))
			}).
			Filter(r.Row.Field("idUser").Eq(idUser)).
			OrderBy(r.Desc("created_at")).
			Merge(func(quote r.Term) interface{}{
				article:=r.Table("article").Get(quote.Field("idArticle"))
				return map[string]interface{}{
					"keywords":article.Field("keywords"),
				}
			}).Run(config.Connection())
		if err!=nil{
			fmt.Println(err)
			return []map[string]interface{}{}
		}
		if result.IsNil(){
			return []map[string]interface{}{}
		}
		var all []interface{}
		err=result.All(&all)
		if err!=nil{
			fmt.Println(err)
			return []map[string]interface{}{}
		}
		quotes:=make([]map[string]interface{},len(all))
		for i,cur:=range all{
			quotes[i]=cur.(map[string]interface{})
		}
		return quotes
	} else{
		return []map[string]interface{}{}
	}
}
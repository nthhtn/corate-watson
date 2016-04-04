package model

import (
	"corate/config"
	r "github.com/dancannon/gorethink"
	"fmt"
)

func InsertArticle(article interface{}) r.WriteResponse{
	result,err:=r.Table("article").Insert(article).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func UpdateArticle(id string,data interface{}) r.WriteResponse{
	result,err:=r.Table("article").Get(id).Update(data).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func GetArticlesByField(field string,value interface{}) []map[string]interface{}{
	result,err:=r.Table("article").Filter(r.Row.Field(field).Eq(value)).
		Run(config.Connection())
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
	articles:=make([]map[string]interface{},len(all))
	for i,cur:=range all{
		articles[i]=cur.(map[string]interface{})
	}
	return articles
}

func GetArticlesByTag(tag string) []map[string]interface{}{
	result,err:=r.Table("article").
		Filter(r.Or(
			r.Row.Field("keywords").Contains(func(keyword r.Term) interface{}{
				return keyword.Field("text").Match("(?i)"+tag)
			}),
			r.Row.Field("taxonomy").Contains(func(taxonomy r.Term) interface{}{
				return taxonomy.Field("text").Match("(?i)"+tag)
			}),
			r.Row.Field("concepts").Contains(func(concept r.Term) interface{}{
				return concept.Field("text").Match("(?i)"+tag)
			}),
		)).Run(config.Connection())
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
	articles:=make([]map[string]interface{},len(all))
	for i,cur:=range all{
		articles[i]=cur.(map[string]interface{})
	}
	return articles
}

func GetKeywords(idArticle string) []interface{}{
	result,err:=r.Table("article").Get(idArticle).Field("keywords").Run(config.Connection())
	defer result.Close()

	if err!=nil{
		fmt.Println(err)
		return []interface{}{}
	}
	var keywords []interface{}
	err=result.All(&keywords)
	if err!=nil{
		fmt.Println(err)
		return []interface{}{}
	}
	return keywords
}
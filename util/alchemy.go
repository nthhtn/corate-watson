package util

import (
	ai "github.com/elvuel/alchemyapi_go"
	"net/url"
	"fmt"
)

var Analyzer *ai.Analyzer

func InitAnalyzer(){
	var err error
	Analyzer,err=ai.NewAnalyzer("d4daae584d558b6082a4f1f1b01972bfbefb86fe")
	if err!=nil{
		panic(err)
	}
}

func Authors(weburl string) []string{
	resp,err:=Analyzer.Authors("url",weburl,url.Values{})
	if err!=nil{
		fmt.Println(err)
		return []string{}
	}
	return resp.Authors.Names
}

func Keywords(weburl string) []map[string]interface{}{
	resp,err:=Analyzer.Keywords("url",weburl,url.Values{})
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	keywords:=make([]map[string]interface{},0)
	for _,k:=range resp.Keywords{
		keywords=append(keywords,map[string]interface{}{
			"text":k.Text,
			"relevance":k.Relevance,
		})
	}
	return keywords
}

func Taxonomy(weburl string) []map[string]interface{}{
	resp,err:=Analyzer.Taxonomy("url",weburl,url.Values{})
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	taxonomy:=make([]map[string]interface{},0)
	for _,t:=range resp.Taxonomies{
		if (t.Score>0.5){
			taxonomy=append(taxonomy,map[string]interface{}{
				"label":t.Label,
				"score":t.Score,
			})
		}
	}
	return taxonomy
}

func Concepts(weburl string) []map[string]interface{}{
	resp,err:=Analyzer.Concepts("url",weburl,url.Values{})
	if err!=nil{
		fmt.Println(err)
		return []map[string]interface{}{}
	}
	concepts:=make([]map[string]interface{},0)
	for _,c:=range resp.Concepts{
		concepts=append(concepts,map[string]interface{}{
			"text":c.Text,
			"relevance":c.Relevance,
		})
	}
	return concepts
}
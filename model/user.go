package model

import (
	"corate/config"
	r "github.com/dancannon/gorethink"
	"fmt"
)

func InsertUser(user interface{}) r.WriteResponse{
	result,err:=r.Table("users").Insert(user).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func UpdateUser(id string,data interface{}) r.WriteResponse{
	result,err:=r.Table("users").Get(id).Update(data).RunWrite(config.Connection())
	if err!=nil{
		fmt.Println(err)
	}
	return result
}

func GetUsersByField(field string,value interface{}) []map[string]interface{}{
	result,err:=r.Table("users").Filter(r.Row.Field(field).Eq(value)).Run(config.Connection())
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
	users:=make([]map[string]interface{},len(all))
	for i,cur:=range all{
		users[i]=cur.(map[string]interface{})
	}
	return users
}
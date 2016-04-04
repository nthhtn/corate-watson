package config

import (
	r "github.com/dancannon/gorethink"
	"fmt"
)

func Connection() *r.Session{
	session,err:=r.Connect(r.ConnectOpts{
		Address:"localhost:28015",
		Database:"corate",
	})
	if err!=nil{
		fmt.Println(err)
	}
	return session
}

func Setupdb(){
	
	// Connect
	session,err:=r.Connect(r.ConnectOpts{
		Address:"localhost:28015",
	})
	if err!=nil{
		fmt.Println(err)
	}

	// Create database
	resp,err:=r.DBCreate("corate").RunWrite(session)
	if err!=nil{
		fmt.Println(err)
	} else{
		fmt.Printf("%d DB created\n",resp.DBsCreated)
	}

	// Create tables
	tables:=[]string{"users","quote","article"}
	for _,table:=range tables{
		resp,err=r.DB("corate").TableCreate(table).RunWrite(session)
		if err!=nil{
			fmt.Println(err)
		} else{
			fmt.Printf("%d table created\n",resp.TablesCreated)
		}
	}

	// Create indexes
	resp,err=r.DB("corate").Table("quote").IndexCreate("idUser").RunWrite(session)
	if err!=nil{
		fmt.Println(err)
	} else{
		fmt.Printf("%d index created\n",resp.Created)
	}
}
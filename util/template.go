package util

import (
    "net/http"
    "html/template"
    "path/filepath"
    "os"
)

func RenderTemplate(w http.ResponseWriter,data interface{},tmpl... string){
    cwd,_:=os.Getwd()
    files:=make([]string, len(tmpl))
    for i,file:=range tmpl{
        files[i]=filepath.Join(cwd,"./view/"+file+".tmpl")
    }
    t,err:=template.ParseFiles(files...)
    if err!=nil{
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }
    templates:=template.Must(t,err)
    err=templates.Execute(w,data)
    if err!=nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
    }
}
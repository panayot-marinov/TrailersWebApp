package src

import (
	"embed"
	"html/template"
)

var tpl *template.Template

//go:embed templates/*
var templatesData embed.FS

func init() {
	//os.Setenv("CONNSTR", "user=postgres password=parolazabaza host=127.0.0.1 port=5432 dbname=MainDB connect_timeout=20 sslmode=disable")
	tpl = template.Must(template.ParseFS(templatesData, "templates/*.html"))
}

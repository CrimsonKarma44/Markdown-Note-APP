package main

import (
	"Markdown_Notetaking_App"
	"log"
	"net/http"
)

func url(mux *http.ServeMux) {
	mux.HandleFunc("/save", Markdown_Notetaking_App.SaveHandler)
	mux.HandleFunc("/list", Markdown_Notetaking_App.ListHandler)
	mux.HandleFunc("/render/", Markdown_Notetaking_App.RenderHandler)
	mux.HandleFunc("/grammar/", Markdown_Notetaking_App.GrammarHandler)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
func main() {
	mux := http.NewServeMux()
	url(mux)
}

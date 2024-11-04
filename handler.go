package Markdown_Notetaking_App

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/save ...")
	var n Note

	file, handler, err := r.FormFile("markdown")
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	n.Title = strings.TrimSuffix(handler.Filename, ".md")
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	n.Body = buf.String()
	n.Save()

	w.Write([]byte("markdown Saved"))
	fmt.Println("/saved ...")
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	var list ListNote
	err := list.get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	value, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(value)
}

func RenderHandler(w http.ResponseWriter, r *http.Request) {
	var n Note
	p := bluemonday.UGCPolicy()
	url := strings.TrimPrefix(r.URL.Path, "/render/")

	n.Get(url)

	v := blackfriday.MarkdownCommon([]byte(n.Body))
	html := p.SanitizeBytes(v)
	w.Write(html)
}

func GrammarHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	var n Note
	url := strings.TrimPrefix(r.URL.Path, "/grammar/")
	n.Get(url)

	api := GrammarApi{
		os.Getenv("API_KEY"),
		n.Body,
		uuid.New(),
	}
	respBody, err := json.Marshal(api)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	resp, err := http.Post(os.Getenv("URL"), "application/json", bytes.NewBuffer(respBody))
	if err != nil {
		fmt.Println("Error making request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	defer resp.Body.Close()

	fmt.Println(os.Getenv("URL"))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	sb := string(body)
	//l, _ := json.Marshal(sb)
	fmt.Println(sb)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(sb))
	//v := json.NewEncoder(w)
	//v.SetIndent("", "  ")
	//v.Encode(sb)

	fmt.Println("response Status:", resp.Status)
}

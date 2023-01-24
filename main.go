package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"log"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

const FILES_TO_SERVE = "/home/plof/Documents/Goproyects/servidorArchivos/config.json"
const TRIES_TO_GET_A_FILE = 10

func everything(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadAll(r.Body)
	if err!=nil{
		  log.Fatal(err)
	}
	fmt.Println(string(input))

	data := gjson.GetBytes(input, "name")
	if data.String() == "neco" {
		file, _ := ioutil.ReadFile("neco-arc.gif")
		fmt.Println("A wild Neco-arc appeared")
		w.Write(file)
		return
	}

	contents, err := ioutil.ReadFile(FILES_TO_SERVE)
	if err!=nil{
		  log.Fatal(err)
	}
	parsed := gjson.ParseBytes(contents)
	value:=parsed.Map()

	for key, element := range value {
		if key == data.String() {
			if element.Get("password").Str == "none" {
				serveFile(gjson.GetBytes(input, "file").String(), element.Get("files"),w)
				return
			} else {
				if element.Get("password").Str == gjson.GetBytes(input, "password").String() {
					serveFile(gjson.GetBytes(input, "file").String(), element.Get("files"),w)
				return
				} else {
					fmt.Println("Wrong password")
					w.Write([]byte("Wrong password"))
					//Add one to attemps getting the file block at specified CONST
					return
				}
			}
		}
		fmt.Println("The name of the share wasnt found")
		fmt.Println(key)
		fmt.Println(data)
		w.Write([]byte("The name of the share wasnt found"))
		return
	}
	fmt.Println("Server Error")
	w.Write([]byte("Server Error"))
}

func serveFile(path string, files gjson.Result,w http.ResponseWriter) {
	if !files.IsArray() {
		fmt.Println("Error in the file configuration, files must be an array")
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("File requested does not exist")
		fmt.Println(path)
		return
	}
	for i := 0; i < len(files.Array()); i++ {
		if files.Get(fmt.Sprintf("%d", i)).String() == path {
			file, err := os.ReadFile(files.Get(fmt.Sprintf("%d",i)).String())
			if err !=nil{
				log.Fatal(err)
			}
			w.Write(file)
			return
		}
	}
	fmt.Println("Cannot serve file")
	w.Write([]byte("Cannot serve file"))
}

func main() {
	fmt.Println("Sever start")
r := mux.NewRouter()
r.HandleFunc("/{_}", everything).Methods("GET")
r.HandleFunc("/", index).Methods("GET")
r.HandleFunc("/", requestHandler).Methods("POST")
http.ListenAndServe(":8000", r)
}

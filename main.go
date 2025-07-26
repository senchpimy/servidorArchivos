package main

import (
	"archive/zip"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

const FILES_TO_SERVE = "/home/plof/Documents/Goproyects/servidorArchivos/config.json"
const TRIES_TO_GET_A_FILE = 10
const PORT = ":8000"
const URL = "http://localhost"

type template_contents struct {
	Port string
	Url  string
}

func everything(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	content := template_contents{Port: PORT, Url: URL}
	err = t.Execute(w, content)
	if err != nil {
		panic(err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	input, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(input))

	data := gjson.GetBytes(input, "name")
	if data.String() == "neco" {
		file, err := os.ReadFile("neco-arc.gif")
		if err != nil {
			log.Println("Error reading neco-arc.gif:", err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		fmt.Println("A wild Neco-arc appeared")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(file)
		return
	}

	contents, err := os.ReadFile(FILES_TO_SERVE)
	if err != nil {
		log.Fatal(err)
	}
	parsed := gjson.ParseBytes(contents)
	value := parsed.Map()

	for key, element := range value {
		if key == data.String() {
			if element.Get("password").Str == "none" {
				verifyFile(gjson.GetBytes(input, "file").String(), element.Get("files"), w)
				return
			} else {
				if element.Get("password").Str == gjson.GetBytes(input, "password").String() {
					verifyFile(gjson.GetBytes(input, "file").String(), element.Get("files"), w)
					return
				} else {
					fmt.Println("Wrong password")
					w.Write([]byte("Wrong password"))
					//Add one to attemps getting the file block at specified CONST
					return
				}
			}
		}
	}

	// Esta es la corrección lógica: solo se ejecuta si el bucle termina sin encontrar una coincidencia.
	fmt.Println("The name of the share wasnt found")
	fmt.Println(data)
	w.Write([]byte("The name of the share wasnt found"))
}

func verifyFile(path string, files gjson.Result, w http.ResponseWriter) {
	if !files.IsArray() {
		fmt.Println("Error in the file configuration, files must be an array")
		w.Write([]byte("Server configuration error"))
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("File requested does not exist:", path)
		w.Write([]byte("File requested does not exist"))
		return
	}
	for _, fileInConfig := range files.Array() {
		if fileInConfig.String() == path {
			serveFile(w, path)
			return
		}
	}
	fmt.Println("Cannot serve file: Not allowed in config")
	w.Write([]byte("Cannot serve file"))
}

func serveFile(w http.ResponseWriter, path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		w.Write([]byte("Error obtaining info about the file"))
		return
	}

	if fileInfo.IsDir() {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", filepath.Base(path)))

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		entries, err := os.ReadDir(path)
		if err != nil {
			log.Printf("Error reading directory %s: %v", path, err)
			http.Error(w, "Could not read directory", http.StatusInternalServerError)
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := filepath.Join(path, entry.Name())
			fileToZip, err := os.Open(filePath)
			if err != nil {
				log.Printf("Failed to open file %s: %v", filePath, err)
				continue
			}
			defer fileToZip.Close()

			zipEntry, err := zipWriter.Create(entry.Name())
			if err != nil {
				log.Printf("Failed to create entry for %s in zip: %v", entry.Name(), err)
				continue
			}

			_, err = io.Copy(zipEntry, fileToZip)
			if err != nil {
				log.Printf("Failed to copy file %s to zip: %v", filePath, err)
				continue
			}
		}

		return

	} else {
		file, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(file)
	}
}

func main() {
	fmt.Println("Server start on port " + PORT)
	r := mux.NewRouter()
	r.HandleFunc("/{_}", everything).Methods("GET")
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/", requestHandler).Methods("POST")
	http.ListenAndServe(PORT, r)
}

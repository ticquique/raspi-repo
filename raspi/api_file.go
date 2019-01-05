/*
 * File share api
 *
 * File share api.
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func GetFileById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	stmt, err := db.Preparex(`SELECT * FROM file WHERE Id=? LIMIT 1`)
	file := File{}

	err = stmt.Get(&file, params["fileId"])

	if err != nil {
		http.Error(w, "No value found with id "+params["fileId"], http.StatusInternalServerError)
		return
	}

	jsonfile, err := json.Marshal(file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonfile)
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Preparex(`SELECT * FROM file`)
	files := []File{}

	err = stmt.Select(&files)

	if err != nil {
		http.Error(w, "No value found with id "+params["fileId"], http.StatusInternalServerError)
		return
	}

	jsonfile, err := json.Marshal(files)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonfile)
}

func NewFile(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")

	switch mimeType {
	case "video/mp4":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/mpeg":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/ogg":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/quicktime":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/webm":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/x-msvideo":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	case "video/x-ms-wmv":
		if e := saveFile(w, file, handle, "film"); e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "No valid mimetype", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader, type_ string) error {
	name := handle.Filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("Couldn't read the file")
	}

	if _, err := os.Stat("./files/" + type_ + "/" + name); !os.IsNotExist(err) {
		return errors.New("File exist")
	}

	e := os.MkdirAll("./files/"+type_, os.ModePerm)
	if e != nil {
		return e
	}

	err = ioutil.WriteFile("./files/"+type_+"/"+name, data, 0666)
	if err != nil {
		return err
	}

	sqlQuery := "INSERT INTO file (Alias, Route, Name, Type) VALUES (?, ?, ?, ?)"
	alias := strings.TrimRight(filepath.Base(name), filepath.Ext(name))
	route := "./files/" + type_ + "/" + name
	db.MustExec(sqlQuery, alias, route, name, type_)
	return nil
}

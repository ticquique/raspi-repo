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
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type FileMetadata struct {
	title        string
	summary      string
	image        *multipart.File
	imageHandler *multipart.FileHeader
}

func GetFileById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	stmt, err := db.Preparex(`SELECT * FROM file WHERE id=? LIMIT 1`)

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

	query := r.URL.Query()
	var stmt *sqlx.Stmt
	var where []string
	values := []interface{}{}

	if query.Get("type") != "" {
		values = append(values, query.Get("type"))
		where = append(where, fmt.Sprintf("%s = ?", "Type"))
	}

	if query.Get("filename") != "" {
		values = append(values, query.Get("filename")+"%")
		where = append(where, fmt.Sprintf("%s LIKE ?", "Filename"))
	}

	if len(query) > 0 {
		stmt, _ = db.Preparex("SELECT * FROM file WHERE " + strings.Join(where, " AND "))
	} else {
		stmt, _ = db.Preparex("SELECT * FROM file")
	}

	files := []File{}

	if err := stmt.Select(&files, values...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	defer file.Close()
	image, imageHandle, imageErr := r.FormFile("image")

	metadata := FileMetadata{
		title:   r.FormValue("title"),
		summary: r.FormValue("summary"),
	}

	if imageErr == nil {
		metadata.image = &image
		metadata.imageHandler = imageHandle
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mimeType := handle.Header.Get("Content-Type")

	switch mimeType {
	case "video/mp4", "video/mpeg", "video/ogg", "video/quicktime", "video/webm", "video/x-msvideo", "video/x-ms-wmv":
		file, e := saveFile(&file, handle, "films", metadata)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		jsonfile, err := json.Marshal(*file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonfile)

	default:
		http.Error(w, "No valid mimetype", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func saveFile(file *multipart.File, handle *multipart.FileHeader, type_ string, metadata FileMetadata) (*File, error) {
	filename := handle.Filename
	data, err := ioutil.ReadAll(*file)
	image := ""
	summary := ""
	title := ""
	assetsDir := os.Getenv("assetsDir")

	switch mimetype := handle.Header.Get("Content-Type"); mimetype {
	case "video/mp4", "video/mpeg", "video/ogg", "video/quicktime", "video/webm", "video/x-msvideo", "video/x-ms-wmv":
	case "image/gif", "image/png", "image/jpeg", "image/bmp", "image/webp", "image/vnd.microsoft.icon", "image/svg+xml", "image/tiff":
	default:
		return nil, errors.New("No valid mimetype")
	}

	if err != nil {
		return nil, errors.New("Couldn't read the file")
	}

	if _, err := os.Stat(fmt.Sprintf("./%s/%s/%s", assetsDir, type_, filename)); !os.IsNotExist(err) {
		return nil, errors.New("File exist")
	}

	if err := os.MkdirAll(fmt.Sprintf("./%s/%s/", assetsDir, type_), os.ModePerm); err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./%s/%s/%s", assetsDir, type_, filename), data, 0666); err != nil {
		return nil, err
	}

	if metadata.image != nil && metadata.imageHandler != nil {
		image = fmt.Sprintf("./%s/%s/%s", assetsDir, "film_assets", metadata.imageHandler.Filename)
		if _, err := saveFile(metadata.image, metadata.imageHandler, "film_assets", FileMetadata{}); err != nil {
			return nil, err
		}
	}

	if metadata.summary != "" {
		summary = metadata.summary
	}
	if metadata.title != "" {
		title = metadata.title
	}

	sqlQuery := "INSERT INTO file (title, route, filename, type, summary, image) VALUES (?, ?, ?, ?, ?, ?)"
	route := fmt.Sprintf("./%s/%s/%s", assetsDir, type_, filename)

	result, err := db.MustExec(sqlQuery, title, route, filename, type_, summary, image).LastInsertId()
	if err != nil {
		return nil, err
	}

	return &File{Title: title, Id: result, Image: image, Filename: filename, Route: route, Summary: summary, Type_: type_}, nil
}

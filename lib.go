package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func conn() (*sql.Conn, error) {
	connString := os.Getenv("MYSQL_CONN_STRING")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func writeData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-type", "application/json")

	result := map[string]interface{}{
		"Status":  http.StatusOK,
		"Data":    data,
		"Message": "",
	}

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("ERROR", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, err error) {
	log.Println("ERROR", err.Error())

	result := map[string]interface{}{
		"Status":  http.StatusInternalServerError,
		"Data":    nil,
		"Message": err.Error(),
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("ERROR", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

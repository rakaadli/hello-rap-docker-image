package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env is required")
	}

	instanceID := os.Getenv("INSTANCE_ID")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "http method not allowed", http.StatusBadRequest)
			return
		}

		text := "hello world"
		if instanceID != "" {
			text = text + ". from " + instanceID
		}

		w.Write([]byte(text))
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllUsersHandler(w, r)
		case "POST":
			createUserHandler(w, r)
		default:
			http.Error(w, "http method not allowed", http.StatusBadRequest)
			return
		}
	})
	server := new(http.Server)
	server.Handler = mux
	server.Addr = "0.0.0.0:" + port

	log.Println("server starting at", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := new(User)
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		writeError(w, err)
		return
	}

	conn, err := conn()
	if err != nil {
		writeError(w, err)
		return
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(context.Background(), "INSERT INTO users(first_name, last_name, birth) VALUES(?, ?, ?)")
	if err != nil {
		writeError(w, err)
		return
	}

	stmtRes, err := stmt.ExecContext(context.Background(), payload.FirstName, payload.LastName, payload.Birth)
	if err != nil {
		writeError(w, err)
		return
	}

	id, _ := stmtRes.LastInsertId()
	result := map[string]interface{}{"LastInsertID": id}

	writeData(w, result)
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := conn()
	if err != nil {
		writeError(w, err)
		return
	}
	defer conn.Close()

	qry, err := conn.QueryContext(context.Background(), "SELECT * FROM users")
	if err != nil {
		writeError(w, err)
		return
	}

	result := make([]User, 0)
	for qry.Next() {
		var id sql.NullInt32
		var firstName sql.NullString
		var lastName sql.NullString
		var birth sql.NullTime
		err = qry.Scan(&id, &firstName, &lastName, &birth)
		if err != nil {
			writeError(w, err)
			return
		}

		user := User{}
		user.ID = int(id.Int32)
		user.FirstName = firstName.String
		user.LastName = lastName.String
		user.Birth = birth.Time

		result = append(result, user)
	}

	writeData(w, result)
}

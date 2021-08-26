package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"time"
)

/**
web 服务端
*/
func main() {
	server := http.Server{
		Addr: "0.0.0.0:9991",
	}

	http.HandleFunc("/", handlerFuture)
	http.HandleFunc("/mysql", getMysqlData)

	server.ListenAndServe()
}

func handlerFuture(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(r.URL.RawQuery)
	fmt.Fprintf(w, "{\n    \"l1\": {\n        \"l1_1\": [\n            \"l1_1_1\",\n            \"l1_1_2\"\n        ],\n        \"l1_2\": {\n            \"l1_2_1\": 121\n        }\n    },\n    \"l2\": {\n        \"l2_1\": null,\n        \"l2_2\": true,\n        \"l2_3\": {}\n    }\n}")
}

type Result struct {
	Code  int
	Param string
	Msg   string
	Data  []User
}

type User struct {
	Id        int
	Username  string
	Password  string
	CreatedAt time.Time
}

func getMysqlData(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql",
		"root:mobikok@2020@(localhost:3306)/young?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT id,username,password,created_at FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	ret := new(Result)
	ret.Code = 200
	ret.Msg = "success"
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
		ret.Data = append(ret.Data, u)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	ret2 := ret

	retJson, err := json.Marshal(ret2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(string(retJson))
	io.WriteString(w, string(retJson))

}

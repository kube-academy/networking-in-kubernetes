package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/go-errors/errors"
)
import _ "github.com/go-sql-driver/mysql"

type frontend struct{}
func(h frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := "<h1>Front-End Home Page</h1><br><ul>"

	dbMsg, _ := extQuery("http://dbsvc.back-end/dbhealth")
	msg = msg+ "<li>Back end DB service reports:" + dbMsg + "<br>"

	extMsg, _ := extQuery("http://extsvc.back-end/externalhealth")
	msg = msg+ "<li>Back end external service reports:" + extMsg + "<br>"

	sourceIP := r.Header.Get("X-FORWARDED-FOR")
	if sourceIP == "" {
		sourceIP = r.RemoteAddr
	}
	msg = msg+ "</ul><p>You are calling from " + sourceIP + "<p>"


	w.Write([]byte(msg))
}

type dbHealth struct{}
func(h dbHealth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := "db connection: ok"
	err := querydb()
	if err != nil {
		msg = "db connection: fail"
	}
	w.Write([]byte(msg))
}

type extHealth struct{}
func(h extHealth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := "external service connection: ok"
	_, err := extQuery(*exturl)
	if err != nil {
		msg = "external service connection: fail"
	}
	w.Write([]byte(msg))

}

var dbserver *string
var exturl *string

func main() {
	dbserver = flag.String("db", "localhost:3306", "MySQL host:port")
	exturl = flag.String("external", "http://www.google.com", "External host URL")
	flag.Parse()

	http.Handle("/externalhealth", extHealth{})
	http.Handle("/dbhealth", dbHealth{})
	http.Handle("/", frontend{})
	log.Fatal(http.ListenAndServe(":9000", nil))

}

func extQuery(url string) (string, error) {
	print("Attempting connect to: "+url+" ")

	r, err := http.Get(url)
	if err != nil {
		println("FAILED")
		println(errors.Wrap(err,1).ErrorStack())

		//log.Fatal(err)
		return "FAILED TO CONNECT", err
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//log.Fatal(err)
		return "FAILED TO CONNECT", err
	}

	println(data)
	return string(data), nil
}


func querydb() error {
	db, err := sql.Open("mysql", "root:example@tcp(" + *dbserver + ")/mysql")

	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer db.Close()


	var (
		host string
		name string
	)
	rows, err := db.Query("select Host, User from user where User = ?", "root")
	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&host, &name)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		log.Println(host, name)
	}
	err = rows.Err()
	if err != nil {
		//log.Fatal(err)
		return err
	}

	return nil
}

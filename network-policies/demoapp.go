package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/go-errors/errors"
	"os"
	"fmt"
	"strings"
	"net"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

var (
	Fail = Red(  "no connection ")
	Pass = Green("connected     ")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}


type hello struct{}
func(h hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println("Hello called")
	msg := Pass+"\n"
	w.Write([]byte(msg))
}

type frontend struct{}
func(h frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println("health called")
	sourceIP := r.Header.Get("X-FORWARDED-FOR")
	if sourceIP == "" {
		sourceIP = r.RemoteAddr
	}
	msg := "Request from: " + sourceIP + "\n"

	hostname, _ := os.Hostname()
	msg = msg + "hostname: "+hostname+ "\n"

	servicePrefixes := map[string]string {"web": "http://web.front-end",
		"dbsvc": "http://dbsvc.back-end",
		"extsvc": "http://extsvc.back-end" }

	for prefix, url := range servicePrefixes {
		// Avoid endless loop, don't call yourself!
		if !strings.HasPrefix(hostname, prefix) {
			msg = msg + "\t (In cluster) " + url + ": " + extQuery(url) + "\n"
		}
	}

	msg = msg + "\t (External)   " + *dbserver + ": " + querydb() + "\n"
	msg = msg + "\t (External)   " + *exturl + ": " + extQuery(*exturl) + "\n"

	w.Write([]byte(msg))
}

var dbserver *string
var exturl *string

func main() {
	dbserver = flag.String("db", "localhost:3306", "MySQL host:port")
	exturl = flag.String("external", "http://www.google.com", "External host URL")
	flag.Parse()

	http.Handle("/",hello{})
	http.Handle("/health", frontend{})
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func extQuery(url string) string {
	print("Attempting connect to: "+url+" ")

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).Dial,
		},
	}

	r, err := client.Get(url)
	if err != nil {
		println("FAILED")
		println(errors.Wrap(err,1).ErrorStack())
		return Fail
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(errors.Wrap(err,1).ErrorStack())
		return Fail
	}

	println(data)
	return Pass
}


func querydb() string {
	db, err := sql.Open("mysql", "root:example@tcp(" + *dbserver + ")/mysql?timeout=3s")

	if err != nil {
		println(errors.Wrap(err,1).ErrorStack())
		return Fail
	}
	defer db.Close()

	var (
		host string
		name string
	)
	rows, err := db.Query("select Host, User from user where User = ?", "root")
	if err != nil {
		println(errors.Wrap(err,1).ErrorStack())
		return Fail
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&host, &name)
		if err != nil {
			println(errors.Wrap(err,1).ErrorStack())
			return Fail
		}
		log.Println(host, name)
	}
	err = rows.Err()
	if err != nil {
		println(errors.Wrap(err,1).ErrorStack())
		return Fail
	}

	return Pass
}

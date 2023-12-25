package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	sqlite "github.com/mattn/go-sqlite3"
)

// a channel as a global variable
var ch chan int64

func validate(y int64) int64 {
	// check if we can issue a HTTP request in the sqlite custom function
	requestURL := "https://example.com"
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	} else {
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("http get request body length:", len(string(body)))
	}
	ch <- y
	return y % 2
}

func main() {

	ch = make(chan int64)

	go func() {
		for {
			select {
			case y := <-ch:
				fmt.Println("go func1 get y:", y)
			}
		}
	}()
	go func() {
		for {
			select {
			case y := <-ch:
				fmt.Println("go func2 get y:", y)
			}
		}
	}()

	sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("validate", validate, false); err != nil {
				return err
			}
			return nil
		},
	})

	db, err := sql.Open("sqlite3_custom", ":memory:")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	var i int64
	err = db.QueryRow("SELECT validate(3)").Scan(&i)
	if err != nil {
		log.Fatal("POW query error:", err)
	}
	fmt.Println("validate(3) =", i) // 0

	err = db.QueryRow("SELECT validate(442)").Scan(&i)
	if err != nil {
		log.Fatal("POW query error:", err)
	}
	fmt.Println("validate(442) =", i) // 0

	_, err = db.Exec("create table foo (department integer, profits integer)")
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// https://stackoverflow.com/q/27214529/864438 sez error-message must be constant
	mkTrigger := `
	CREATE TRIGGER insert_trigger
	BEFORE INSERT ON foo
	WHEN validate(NEW.department) != 0
	BEGIN
	  SELECT RAISE(ABORT, 'bad validate');
	END;
	`
	_, err = db.Exec(mkTrigger)
	if err != nil {
		log.Fatalf("error creating trigger %q: %s", err, mkTrigger)
		return
	}

	log.Println("inserting 1st")
	_, err = db.Exec("insert into foo values (10, 10)")
	if err != nil {
		log.Println("Failed to insert first:", err)
	}

	log.Println("inserting 2nd")
	_, err = db.Exec("insert into foo values (11, 10)")
	if err != nil {
		log.Println("Failed to insert 2nd:", err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	sqlite "github.com/mattn/go-sqlite3"
	"log"
)

func validate(y int64) int64 {
	return y % 2
}

func main() {
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
package main

import "database/sql"

var db, err = sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")

func dbSettings() {
	if err != nil {
		panic(err)
	}
}

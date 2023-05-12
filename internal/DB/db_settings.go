package DB

import "database/sql"

var Db, err = sql.Open("sqlserver", "server=localhost;user Id=Artem;password=sql12345678;database=Students;encrypt=disable")

func DbSettings() {
	if err != nil {
		panic(err)
	}
}

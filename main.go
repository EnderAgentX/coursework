package main

import (
	"os"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	w := App()
	dbSettings()
	ReadStudents(db)
	ReadGroup(db)
	w.ShowAndRun()
	defer db.Close()
}

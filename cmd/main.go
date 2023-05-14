package main

import (
	"github.com/EnderAgentX/coursework/internal/App"
	"github.com/EnderAgentX/coursework/internal/DB"
	"os"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	w := App.App()
	DB.DbSettings()
	DB.ReadStudents(DB.Db)
	DB.ReadGroup(DB.Db)
	w.ShowAndRun()
	defer DB.Db.Close()
}

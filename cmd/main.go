package main

import (
	"github.com/EnderAgentX/coursework/internal/DB"
	"github.com/EnderAgentX/coursework/internal/app"
	"os"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	w := app.App()
	DB.DbSettings()
	DB.ReadStudents(DB.Db)
	DB.ReadGroup(DB.Db)
	w.ShowAndRun()
	defer DB.Db.Close()
}

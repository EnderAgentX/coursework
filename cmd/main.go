package main

import (
	"github.com/EnderAgentX/coursework/internal"
	"os"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	w := internal.App()
	internal.DbSettings()
	internal.ReadStudents(internal.Db)
	internal.ReadGroup(internal.Db)
	w.ShowAndRun()
	defer internal.Db.Close()
}

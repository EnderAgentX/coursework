package model

type Student struct {
	Id          int
	Name        string
	Gender      string
	StudentCard string
	Phone       string
	GroupId     int
}

type Group struct {
	Id   int
	Name string
}

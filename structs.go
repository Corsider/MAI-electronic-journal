package main

//MAIUSER
type MAIUSER struct {
	//student
	Name        string `bson:"name"`
	Disciplines string `bson:"disciplines"`
	Email       string `bson:"email"`
	Marks       string `bson:"marks"`
	Token       string `bson:"Token"`
	//teacher
	Initials   string `bson:"initials"`
	Discipline string `bson:"discipline"`
	Isteacher  bool   `bson:isteacher`
}

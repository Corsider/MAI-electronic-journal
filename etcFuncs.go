package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

//GetNews ...
func GetNews() (bool, string) {
	stuffCollection := mpDatabaseRS.Collection("stuff")
	//filter, err := stuffCollection.Find(ctx, bson.M{"stuff_type": 0})
	//if err != nil {
	//	log.Fatal(err)
	//}
	var news bson.M
	//if err = filter.Decode(&news); err != nil {
	//	log.Fatal(err)
	//}
	err := stuffCollection.FindOne(ctx, bson.M{"stuff_type": 0}).Decode(&news)
	if err != nil {
		log.Println(err)
		//log.Fatal(err)
	}
	if news["is_news_available"] == false {
		return false, ""
	} else {
		newsString := news["current_news"].(string)
		return true, newsString
	}
}

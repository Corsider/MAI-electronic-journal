package main

import "github.com/gin-gonic/gin"

//InitRouters ...
func InitRouters(r *gin.Engine) {
	//General section
	r.GET("/ping", UPing)                   //checks connection
	r.GET("/createAccount", UCreateAccount) //Creates account

	//Student section
	r.GET("/action/s/", isAuthorized(ActionStudent, "s"))

	//Teacher section
	r.GET("/action/t/", isAuthorized(ActionTeacher, "t"))
}

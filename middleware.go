package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var mySigningKey = []byte("lab11wasHARD123")

// GenerateJWT generates JWT token
func GenerateJWT(mpJWTuser string, accType string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = mpJWTuser
	claims["utype"] = accType

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Something wrong while generating JWT token: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

//isAuthorized middleware protects endpoints from unauthorized access
func isAuthorized(endpoint func(c *gin.Context), userType string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(string(c.GetHeader("Token")), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("[SERVER] Error while authorization process")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Println("[SERVER] Error while auth process.", err)
				c.JSON(401, gin.H{
					"msg": "Authorization fail. Maybe you've entered invalid token?",
				})
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok && token.Valid {
				if claims["utype"] == "yes" && userType == "t" {
					fmt.Println("[SERVER] Teacher performed an action")
					endpoint(c)
				} else if claims["utype"] == "yes" && userType == "s" {
					fmt.Println("[SERVER] Teacher tried to access student's actions")
					c.JSON(200, gin.H{
						"msg": "You tried to access student's actions via teacher's account. Access denied.",
					})
				} else if claims["utype"] == "no" && userType == "t" {
					fmt.Println("[SERVER] Somebody (student may be) tried to access teacher's actions on server! Shame on him!")
					c.JSON(200, gin.H{
						"msg": "You don't have access to this page. Please check if you included auth token (and it is correct) in header of the request. Access denied.",
					})
				} else if claims["utype"] == "no" && userType == "s" {
					fmt.Println("[SERVER] Student performed an action")
					endpoint(c)
				} else {
					c.JSON(200, gin.H{
						"msg": "It seems that your account type is incorrect. Please recreate account and fill is_teacher field with no or yes.",
					})
				}

			}
		} else {
			fmt.Println("[SERVER] User provided no token in the header of request.")
			c.JSON(200, gin.H{
				"msg": "You don't have access to this page. Please check if you included auth token (and it is correct) in header of the request.",
			})
		}
	})
}

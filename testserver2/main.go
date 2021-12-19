package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func sayHello(c *gin.Context){

	issuer := c.Request.Header["Iss"]
	log.Println(issuer)
	log.Println(c.Request.URL.Path)
	c.JSON(200, gin.H{
		"message": c.Param("name"),
	})
}


func main(){

	router := gin.Default()
	router.GET("/c1/p1/:name",sayHello)
	router.Run(":8081")

	
}
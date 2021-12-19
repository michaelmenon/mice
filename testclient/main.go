package main

import (
	
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	
	"os"
	"github.com/golang-jwt/jwt"
	"github.com/gin-gonic/gin"
	"crypto/tls"
)

func sayHello(c *gin.Context){

	log.Println(c.Request.URL.Path)
	c.JSON(200, gin.H{
		"hello from c2": c.Param("name"),
	})
}

func createToken()(string,error){
	key := os.Getenv("mysecret")
	mySigningKey := []byte(key)

	// Create the Claims
	claims := &jwt.StandardClaims{
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss,err
}

func main(){

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	c := &http.Client{Transport: tr}
    req, err := http.NewRequest("GET", "https://localhost:8000/c1/p1/sam", nil)
    if err != nil {
        fmt.Printf("error %s\n", err)
        return
    }
	/*token,err := createToken()
	if err != nil {
        fmt.Printf("error %s\n", err)
        return
    }
    req.Header.Add("Authorization", `Bearer `+token)*/
    resp, err := c.Do(req)
    if err != nil {
        fmt.Printf("error %s\n", err)
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    fmt.Printf("Body : %s\n", body)
}
	

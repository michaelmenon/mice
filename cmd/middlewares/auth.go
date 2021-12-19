package middlewares

import (
	"os"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Claims map[string]string

func validateJWT(tokenString string)(Claims,error){
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := os.Getenv(viper.GetString("config.authenv"))
		return []byte(secret), nil
	})
	if err!=nil{
		return nil,err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c := make(Claims)	
		for k,v := range claims{
			if value,ok := v.(string);ok{
				c[k] = value
			}
		}
		return c,nil
	} 


	return nil,err

}
//do auth if its set in the toml file
func DoAuth(c *gin.Context)(Claims,error){

	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	if token == ""{
		return nil,fmt.Errorf("Authorization Bearer token not found")
	}
	return validateJWT(token)
}
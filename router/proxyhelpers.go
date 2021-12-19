package router

import (
	"net/url"
	"os"
	"sync/atomic"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/ratelimit"
)
type server struct{
	u *url.URL
	up bool
}
type proxy struct {
	remote []*server
	index uint64
	
}

type Claims map[string]string

var (
	limit ratelimit.Limiter
	proxies map[string]*proxy
	
)

//roundrobin ... do round robil load balancing strategy
func (p *proxy)roundrobin() *url.URL{
	var curr int
	//get the current index but limit it to the number remote server we need to load balance
	curr = int(atomic.AddUint64(&p.index, 1))%len(p.remote)
	l := len(p.remote)
	for i :=curr ; i < l+curr; i++{
		
		id := i%len(p.remote)
		if p.remote[id].up{
		
		//now store the current inex back as the proxy inex atomically
			atomic.StoreUint64(&p.index,uint64(id))
			return p.remote[id].u
		}
	}
	return nil
}

//do rate limiting
func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		prev = now
	}
}

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
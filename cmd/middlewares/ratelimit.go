package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"
	"go.uber.org/ratelimit"
	"github.com/spf13/viper"
)

//do rate limiting
func leakBucket(limit ratelimit.Limiter) gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		prev = now
	}
}

func DoRateLimit(r *gin.Engine){
	//set the ratelimiter if the ratelimit is set in the toml
	if viper.GetBool("config.ratelimit"){
		limit := ratelimit.New(viper.GetInt("config.ratecount"))
		r.Use(leakBucket(limit))
	}
}
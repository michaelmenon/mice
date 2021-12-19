package router

import (
	"fmt"
	

	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/ratelimit"
)

//RegisterProxies ... register the proxies
func RegisterProxies(r *gin.Engine,configfilepath string)error{
	viper.SetConfigFile("mice.toml")
	viper.AddConfigPath(configfilepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	//now load the handlers
	loadHandlers(r)
	//run the health checker
	go runHealthCheck()
	return nil
}

//loadHandlers ... sets all the required handlers for Mice
func loadHandlers(r *gin.Engine){
	
	//set the ratelimiter if the ratelimit is set in the toml
	if viper.GetBool("config.ratelimit"){
		limit = ratelimit.New(viper.GetInt("config.ratecount"))
		r.Use(leakBucket())
	}

	proxies = make(map[string]*proxy)
	
	setProxy()

	for path := range proxies{ 
		
		r.Any(path+"/*any",runProxy)
	}
}
//do auth if its set in the toml file
func doAuth(c *gin.Context)(Claims,error){

	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	if token == ""{
		return nil,fmt.Errorf("Authorization Bearer token not found")
	}
	return validateJWT(token)
}
//doProxy ... set the proxy handlers
func setProxy(){
	servers := viper.GetStringMapString("server")
	for k,_ := range servers{
		k = fmt.Sprintf("server.%s",k)
		addrs := viper.GetStringSlice(k+".addr")
		path := viper.GetString(k+".role")
		var remote []*server
		for _,addr := range addrs{ 
			rm, err := url.Parse(addr)
			if err != nil {
				fmt.Printf("error parsing remote server %s with error %v",addr,err)
				continue
			}
			remote = append(remote,&server{u:rm,up: true})
			
		}
		proxies[path] = &proxy{
			remote:  remote,
			index: 0,
		}
	}
	
}
//run the proxy handler... route to the proxy handlers
func runProxy(c *gin.Context){

	var claims Claims
	var err error
	//check if we need to do auth
	if viper.GetBool("config.doauth"){

		claims,err = doAuth(c)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
	path := c.Request.URL.Path
	paths := strings.Split(path,"/")
	if len(paths) <=1{
		c.JSON(http.StatusBadRequest, gin.H{"error": "no path found"})
		return
	}
	path = "/"+paths[1]
	
	if server,ok := proxies[path];ok{

		//do load balancing
		remote := server.roundrobin()
		if remote == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "all the servers are down"})
			return
		}
		remoteServer := httputil.NewSingleHostReverseProxy(remote)
		remoteServer.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			//add all the claims we got from the token
			//as header values
			//claims keys are canonicalized
			//only Standard claim keys are supported
			for k,v := range claims{
				req.Header.Add(k, v)
			}
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
		}
		remoteServer.ServeHTTP(c.Writer, c.Request)
	}
}
	

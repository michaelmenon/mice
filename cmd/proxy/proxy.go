package proxy

import (
	"fmt"
	
	"mice/cmd/middlewares"

	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	
)

type server struct{
	u *url.URL
	up bool
}
type proxy struct {
	remote []*server
	index uint64
	
}

var (
	proxies map[string]*proxy
)

//RegisterProxies ... register the proxies
func RegisterProxies(r *gin.Engine,configfilepath string)error{
	viper.SetConfigFile("mice.toml")
	viper.AddConfigPath(configfilepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	//now load the handlers
	loadMiddlewares(r)
	//run the health checker
	go runHealthCheck()
	return nil
}

//loadHandlers ... sets all the required handlers for Mice
func loadMiddlewares(r *gin.Engine){
	
	middlewares.DoRateLimit(r)
	proxies = make(map[string]*proxy)
	
	setProxy()

	for path := range proxies{ 
		if !strings.HasPrefix(path,"/"){
			fmt.Printf("Mice:[ERR]","Path %s does not begin with forward slash, so will not be added for routing\n",path)
			continue
		}
		if strings.Count(path,"/") > 1{
			fmt.Printf("Mice:[ERR]","Path %s has more than one forward slash, so will not be added for routing\n",path)
			continue
		}
		r.Any(path+"/*any",runProxy)
	}
	
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

	var claims middlewares.Claims
	var err error
	//check if we need to do auth
	if viper.GetBool("config.doauth"){

		claims,err = middlewares.DoAuth(c)
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
		remote := roundrobin(server)
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
	

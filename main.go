package main

import (
	"fmt"
	"os"
	"mice/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main(){

	gin.SetMode(gin.ReleaseMode)
    gr := gin.Default()

	//register proxies from the TOML file
	err := router.RegisterProxies(gr)
	if err!=nil{
		_ = fmt.Errorf("error registering the proxies from toml: %v\n",err)
		os.Exit(1)
	}
	
	
	server := viper.GetStringMap("gateway")
	addr := fmt.Sprintf("%s:%s",server["ip"],server["port"])
	
	if viper.GetBool("config.tls"){
		fmt.Printf("Running the Mice Gateway over TLS on => %s\n",addr)
		gr.RunTLS(addr,viper.GetString("config.cert"),viper.GetString("config.key"))
	}else{
		fmt.Printf("Running the Mice Gateway on => %s\n",addr)
		gr.Run(addr)
	}

}


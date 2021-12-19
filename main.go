package main

import (
	"fmt"
	"os"
	"flag"
	"mice/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main(){

	config := flag.String("config","","toml file path")
	flag.Parse()

	if *config == ""{
		fmt.Println("Mice [ERR]:No toml file path specified. Run mice -h for details.",)
		os.Exit(1)
	}
	gin.SetMode(gin.ReleaseMode)
    gr := gin.Default()

	//register proxies from the TOML file
	err := router.RegisterProxies(gr,*config)
	if err!=nil{
		fmt.Printf("Mice [ERR]:Error registering the proxies from toml: %v\n",err)
		os.Exit(1)
	}
	
	
	server := viper.GetStringMap("gateway")
	addr := fmt.Sprintf("%s:%s",server["ip"],server["port"])
	
	if viper.GetBool("config.tls"){
		fmt.Printf("Mice [INFO]:Running the Mice Gateway over TLS on => %s\n",addr)
		gr.RunTLS(addr,viper.GetString("config.cert"),viper.GetString("config.key"))
	}else{
		fmt.Printf("Mice [INFO]:Running the Mice Gateway on => %s\n",addr)
		gr.Run(addr)
	}

}


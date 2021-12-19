package server

import (
	"flag"
	"fmt"
	"log"
	"mice/cmd/constants"
	"mice/cmd/proxy"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaelmenon/mice/cmd/constants"
	"github.com/spf13/viper"
)

func Run(){
	config := flag.String("config","","toml file path without the file name. if the file is in the current folder the path as ./")
	flag.Parse()

	if *config == ""{
		fmt.Printf("%sNo toml file path specified. Run mice -h for details.",constants.MICEERR)
		os.Exit(1)
	}
	gin.SetMode(gin.ReleaseMode)
    gr := gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	gr.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s[%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				constants.MICEINFO,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	gr.Use(gin.Recovery())
	//register proxies from the TOML file
	err := proxy.RegisterProxies(gr,*config)
	if err!=nil{
		fmt.Printf("Mice -[ERR]::Error registering the proxies from toml: %v\n",err)
		os.Exit(1)
	}
	
	
	server := viper.GetStringMap("gateway")
	addr := fmt.Sprintf("%s:%s",server["ip"],server["port"])
	
	if viper.GetBool("config.tls"){

		fmt.Printf("%sRunning the Mice Gateway over TLS on => %s\n",constants.MICEINFO,addr)
		log.Fatal(gr.RunTLS(addr,viper.GetString("config.cert"),viper.GetString("config.key")))
	}else{
		fmt.Printf("%sRunning the Mice Gateway on => %s\n",constants.MICEINFO,addr)
		log.Fatal(gr.Run(addr))
	}
}

package proxy

import (
	"fmt"
	"net"
	"time"
	"mice/cmd/constants"
)

//RunHealthCheck start health checking routine for the registered proxies
func runHealthCheck(){

	ticker := time.NewTicker(5 *time.Second)
	defer func(){
		ticker.Stop()
	}()
	for {
		select{
		case <-ticker.C:
			dialServer()
		}
	}
}

//dialServer .. dial the proxies to check their health
func dialServer(){
	for _,v := range proxies{
		for _,u := range v.remote{

			ip := u.u.Host
			c,err := net.Dial("tcp",ip,)
			if err!=nil{
				fmt.Printf("%sserver with address %s is down\n",constants.MICEERR,ip)
				u.up = false
			}else{
				u.up = true
				c.Close()
			}
			
		}
	}
}
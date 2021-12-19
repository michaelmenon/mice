package proxy

import (

	"sync/atomic"

	"net/url"
)

//it only supports round robin for now.. plan is to support other
//strategies later

//roundrobin ... do round robil load balancing strategy
func roundrobin(p *proxy) *url.URL{
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
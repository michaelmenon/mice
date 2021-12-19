# mice
A simple API Gateway in Golang.



<img width="277" alt="Screen Shot 2021-12-18 at 11 57 43 PM" src="https://user-images.githubusercontent.com/5271064/146664415-d3de7881-848e-4bcc-84aa-e5810319062f.png">




Installation:
Using go install:
  1) First install Go, if not already there
  2) Set GOPATH and GOBIN if not already set as in step #3, if already set skipe to step#8
  3) mkdir $HOME/go
  4) mkdir $HOME/go/bin
  5) export GOPATH=$HOME/go
  7) export GOPATH=$HOME/go/bin
  8) go get github.com/michaelmenon/mice
  9) go install github.com/michaelmenon/mice@latest

Copy the sample TOML config file provided in the github repo to the path from you are going to run the mice gasteway. Open the terminal and go the the path where you have copied thr toml file and run the mice with this command : mice -config=./

If you want Gateway to run over tls,place the tls certificate and key to the same folder where the toml file is and you need to provide the file path for the cert and key in the toml.

It has the following features:

1) Does reverse proxy to the servers provided in the TOML file:
  For eg: 
  
  [server.1]
  role="/c1"
  addr=["http://localhost:8080","http://localhost:8081"]
  
  [server.2]
  role="/c2"
  addr=["http://localhost:8083","http://localhost:8084"]
  
  Here we have defined 2 servers to be proxied based on the path mentioned in the "role"
  Each server has multiple  instances as mentioned in the array "addr". Mice will do round robin load balancing to the multiple instances in the array.
  
2) Does load balancing to the server if the server has multiple instances provided in the "addr" tag.
3) Does health check for each server instance and log if any server is down
4) Does logging 
5) Supports Authentication Bearer Token if the "doauth" flag is set to true. If its set to true then it will check for the JWT secret key in the env variable set with the name in the tag "authenv". 
If there are any claims Mice will collect that and pass it in the request headers to the corresponding proxies.
Only standard claims are supported as give below :
<img width="397" alt="Screen Shot 2021-12-19 at 9 37 12 AM" src="https://user-images.githubusercontent.com/5271064/146678829-4df4e3cd-ed98-497d-b3f9-676d00ce5c38.png">

7) Does rate limiting if the tag "ratelimit" is set to true. Mention the rate limit count in the tag "ratecount"
8) Support TLS if you set the flag "tls" to true and set the "cert" key to the certificate file path and "key" tag to the key file path.


Sample TOML file : 
<img width="682" alt="Screen Shot 2021-12-18 at 11 54 35 PM" src="https://user-images.githubusercontent.com/5271064/146664378-2c70fd31-f552-4da0-a4e2-55b46caa588b.png">




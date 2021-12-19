# mice
A simple API Gateway in Golang.

Installation:
1) First install Go
2) go get github.com/michaelmenon/mice
3) go install github.com/michaelmenon/mice

Copy the sample TOML config file provided in the github repo to the path from you are going to run the mice gasteway. Open the terminal and go the the path where you have copied thr toml file and run the mice with this command : mice -config=./

Place the tls certificate and key to the same folder where the toml file is and you need to provide the file path for the cert and key in the toml.

It has the following features:

1) Do reverse proxy to the servers provided in the TOML file:
  For eg: 
  
  [server.1]
  role="/c1"
  addr=["http://localhost:8080","http://localhost:8081"]
  
  [server.2]
  role="/c2"
  addr=["http://localhost:8083","http://localhost:8084"]
  
  Here we have defined 2 servers to be proxied based on the path mentioned in the "role"
  Each server has multiple  instances as mentioned in the array "addr". Mice will do round robin load balancing to the multiple instances in the array.
  
2) Do load balancing to the server if the server has multiple instances provided in the "addr" tag.
3) Do health check for each server instance and log if any server is down
4) Do logging 
5) Support Authentication Berer Tomen if the "doauth" flag is set to true. If its set to true then it will check for the JWT secret key in the env variable set in the tag "authenv". If there are any claims Mice will collect that and pass it in the request headers to the corresponding proxies.
6) Do rate limiting if the tag "ratelimit" is set to true. Mention the rate limit count in the tag "ratecount"
7) Support TLS if you set the flag "tls" to true and set the "cert" key to the certificate file path and "key" tag to the key file path.


Sample TOML file : 

<img width="532" alt="Screen Shot 2021-12-18 at 11 42 47 PM" src="https://user-images.githubusercontent.com/5271064/146664150-214ceca6-e55d-48f3-9b7d-7074b040c355.png">


# Go-Mongo-REST-API
## What is this...
This is an API that serves products with applied prices. 
- It is written in Go.
- It uses Fiber an excellent web framework built on top of as excellent Fasthttp and expired by also as excellent ExpressJS.
- It uses MongoDB as data storage.

It is as simple as they can. There's a separation between Business Logic and the rest of the application (adapters). 
Said adapters make use of Business Logic rather than otherwise in accordance to DDD and the hexagonal paradigm.
Business logic is contained inside `internal/model` package.
While the idea is to deploy this using `docker-compose`,  
**Kubernetes** deployment / **Docker Swarm** stack can be used by adding/changing a few .ymls 

## Running this
The preferred way to run this would be:  
>$ **make dockerup**

(see *Makefile* for more build targets.)

Once both services are up, API will be available at `http://localhost:8080` 
via `GET /products`
and can be filtered by `priceLessThan=XXX` query param.
  
## Running API documentation
>$ **make swagger** 
> 
and then access `http://localhost:8008`

## What is needed to scale this but is missing for simplicity sake    
This application has been designed to be easily scalable yet doesn't implement everything that is needed any of the 
optional components such as reverse proxy, a key-value storage and a  
message queue or an event bus.
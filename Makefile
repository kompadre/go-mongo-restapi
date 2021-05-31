MYUID = $(shell id -u)
MYGID = $(shell id -g)

.PHONY: dockerup dockerdown

all: bin/server .env assets/extra.json

bin/server: cmd/server/main.go
	go build -o bin/server cmd/server/main.go

.env:
	bash -c "MYUID=$(MYUID) MYGID=$(MYGID) envsubst < .env.example" > .env

dockerup: .env assets/extra.json
	docker-compose up -d --build

dockerdown: .env
	docker-compose down

swagger: openapi.yml
	@echo "Starting swagger editor at http://localhost:8008"
	docker run -p 8008:8080 -v `pwd`:/tmp -e SWAGGER_FILE=/tmp/openapi.yml swaggerapi/swagger-editor


assets/extra.json:
	@wget -q http://dre.red/my/json/extra.json -O assets/extra.json; exit 0

clean:
	rm bin/*
	rm .env
	rm go.sum


MYUID = $(shell id -u)
MYGID = $(shell id -g)

.PHONY: dockerup dockerdown

all: bin/server .env

bin/server: cmd/server/main.go
	go build -o bin/server cmd/server/main.go

.env:
	bash -c "MYUID=$(MYUID) MYGID=$(MYGID) envsubst < .env.example" > .env

dockerup: .env
	docker-compose up -d --build

dockerdown: .env
	docker-compose down

clean:
	rm bin/*
	rm .env


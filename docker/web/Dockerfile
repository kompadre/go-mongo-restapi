FROM go-mongo-restapi-base:staging
WORKDIR /usr/src/app
COPY cmd /usr/src/app/cmd
COPY assets /usr/src/app/assets
COPY internal /usr/src/app/internal
# RUN go mod download
RUN go build cmd/main.go
# EXPOSE 8080
CMD [ "go", "run", "cmd/main.go" ]
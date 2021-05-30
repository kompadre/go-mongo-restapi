FROM go-mongo-restapi-base:staging
WORKDIR /usr/src/app
COPY ./cmd /usr/src/app/cmd
COPY ./assets /usr/src/app/assets
COPY ./internal /usr/src/app/internal
RUN go build -o /server ./cmd/server/main.go
# RUN ./bin/server
# EXPOSE 8080
CMD [ "/server" ]
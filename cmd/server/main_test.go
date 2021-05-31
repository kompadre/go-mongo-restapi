package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerRun(t *testing.T) {
	// Behold something never seen before
	// Goroutines in testing!
	server := Server{}
	os.Setenv("TESTING", "F")
	cwd, _ := os.Getwd()
	os.Chdir( cwd + `/../../` )
	go func() {
		server.run()
	}()

	fmt.Println("Sleeping...")
	time.Sleep(1 * time.Second)
	fmt.Println("Awake...")
	client := http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
		Timeout:       1 * time.Second,
	}
	res, err := client.Get("http://localhost:8080/products")
	if err != nil {
		t.Errorf("Error getting to http server: %v", err)
	}
	assert.Equal(t, res.Status, "200 OK", "http status code should be 200")

	fmt.Println("Shutting down...")
	server.app.Shutdown()

}

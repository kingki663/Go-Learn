package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	fmt.Printf("%s", respBytes)
}

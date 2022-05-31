package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/vmihailenco/msgpack"

	"github.com/the-code-innovator/squeezer/shortener"
)

func port() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return port
}

func main() {
	address := fmt.Sprintf("http://localhost:%s", port())
	shortLink := shortener.ShortLink{}
	shortLink.URL = "https://github.com/the-code-innovator?tab=repositories"

	body, err := msgpack.Marshal(&shortLink)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	msgpack.Unmarshal(body, &shortLink)

	log.Printf("%v\n", shortLink)
}

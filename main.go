package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func ping(url string, respCh chan int, errCh chan error) {
	resp, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	respCh <- resp.StatusCode
}

func main() {
	path := flag.String("file", "urls.txt", "path to URL file")
	flag.Parse()
	file, err := os.ReadFile(*path)
	if err != nil {
		fmt.Println("Ошибка чтения файла")
		panic(err.Error())
	}
	urls := strings.Split(string(file), "\n")

	var validUrls []string
	for _, v := range urls {
		if u := strings.TrimSpace(v); u != "" {
			validUrls = append(validUrls, u)
		}
	}

	var wg sync.WaitGroup
	respCh := make(chan int)
	errCh := make(chan error)
	for _, value := range validUrls {
		wg.Go(func() {
			ping(value, respCh, errCh)
		})
	}
	for range len(validUrls) {
		select {
		case err := <-errCh:
			fmt.Println(err.Error())
		case resp := <-respCh:
			fmt.Printf("Status code: %d\n", resp)
		}
	}
	func() {
		defer close(respCh)
		defer close(errCh)
		wg.Wait()
	}()
}

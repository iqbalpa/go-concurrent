package trash

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type Result struct {
	Url string
	Title string
}

// Correct implementation using goroutine, channel, mutex, defer
func Trash() {
	fmt.Println("================== trash.go ====================")

	// var mu sync.Mutex
	var wg sync.WaitGroup

	urls := []string{"https://github.com", "https://go.dev", "https://gobyexample.com"}
	results := make(map[string]string)
	c := make(chan Result, 3)

	// assign to all goroutines first
	for _,url := range urls {
		wg.Add(1)
		go retrieve(url, c, &wg)
	}

	// This goroutine waits for all worker goroutines to finish (using wg.Wait())
	// and then closes the channel to signal that no more data will be sent.
	go func() {
		wg.Wait() // wait all goroutines finished
		close(c)	// close the channel
	}()

	// No need for a mutex because only the main goroutine accesses the results map.
	// Other goroutines (retrieve) are only sending data via the channel, not modifying the map.
	for res := range c {
		results[res.Title] = res.Url
	}

	fmt.Println(results)
}


func retrieve(url string, c chan Result, wg *sync.WaitGroup) {
	// Defer wg.Done() to ensure that the WaitGroup counter is decremented
	// once the goroutine finishes its execution.
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	res := extractTitle(string(body))

	c <- Result{Url: url, Title: res}
}


func extractTitle(body string) string {
	r, _ := regexp.Compile("<title>(.*?)</title>")
	return r.FindStringSubmatch(body)[1]
}

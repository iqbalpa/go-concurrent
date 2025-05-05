package dump

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

// Correct implementation using goroutine, mutex, defer
func Dump() {
	fmt.Println("================== dump.go ====================")
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// create much urls to trigger race conditions (if not using lock-unlock)
	urls := make([]string, 5000)
	for i := 0; i < 5000; i++ {
			urls[i] = "https://go.dev"
	}

	// The `results` map is shared across all goroutines, so we need to synchronize access using a mutex.
	results := make(map[string]string)

	for _,url := range urls {
		wg.Add(1)
		go retrieve(url, &wg, &mu, results)
	}

	wg.Wait()

	fmt.Println(results)
}


// Locking the mutex to prevent race conditions when accessing the results map.
// This ensures that only one goroutine can write to the map at a time.
func retrieve(url string, wg *sync.WaitGroup, mu *sync.Mutex, results map[string]string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	res := extractTitle(string(body))

	// The `Lock()` ensures that only one goroutine can update the `results` map at a time.
	// Without it, multiple goroutines could modify the map simultaneously, leading to race conditions.
	mu.Lock()
	results[url] = res
	mu.Unlock()
}


func extractTitle(body string) string {
	r, _ := regexp.Compile("<title>(.*?)</title>")
	res := r.FindStringSubmatch(body)
	if len(res) < 2 {
		return ""
	}
	return res[1]
}

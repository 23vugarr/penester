package pkg

import (
	"fmt"
	"net/http"
	"sync"
)

func DirScaner(website string) {
	directories := []string{"admin", "login", "test", "backup", "tmp", "images"}

	var wg sync.WaitGroup
	for _, dir := range directories {
		wg.Add(1)
		go func(directory string) {
			defer wg.Done()
			url := fmt.Sprintf("%s%s/", website, directory)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Failed to reach %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				fmt.Printf("Directory found: %s\n", url)
			} else {
				fmt.Printf("Directory not found: %s (Status code: %d)\n", url, resp.StatusCode)
			}
		}(dir)
	}
	wg.Wait()
}

package pkg

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func PortScanner(start, end int, website string) {
	var wg sync.WaitGroup
	wg.Wait()
	for port := start; port <= end; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", website, port)
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				fmt.Printf("Port %d closed\n", port)
				return
			}
			conn.Close()
			fmt.Printf("Port %d open\n", port)
		}(port)
	}
}

package pkg

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func PortScanner(start, end int, website string, doneCh chan<- bool) {
	var wg sync.WaitGroup
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
			err = conn.Close()
			if err != nil {
				return
			}
			fmt.Printf("Port %d open\n", port)
		}(port)
	}
	wg.Wait()
	log.Println("sending done...")
	doneCh <- true
	close(doneCh)
}

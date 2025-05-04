package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func portScan(host string, port int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			results <- fmt.Sprintf("Port %d: Closed or filtered\n", port)
		} else {
			results <- fmt.Sprintf("Port %d: Error: %s\n", port, err)
		}
		return
	}
	defer conn.Close()
	results <- fmt.Sprintf("Port %d: Open\n", port)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)

	if len(os.Args) < 3 {
		fmt.Println("Usage: port_scanner <target_host> <start_port> [end_port]")
		return
	}

	targetHost := os.Args[1]
	startPort, _ := strconv.Atoi(os.Args[2])
	endPort := startPort

	if len(os.Args) == 4 {
		endPort, _ = strconv.Atoi(os.Args[3])
	}

	fmt.Printf("Scanning ports on %s from %d to %d...\n", targetHost, startPort, endPort)

	for port := startPort; port < endPort; port++ {
		wg.Add(1)
		go portScan(targetHost, port, ch, &wg)
	}
	// Запускаем горутину для чтения из канала
	go func() {
		for msg := range ch {
			fmt.Println(msg)
		}
	}()

	wg.Wait() // Ждем завершения всех горутин, включая горутину, которая читает из канала
	close(ch) // Закрываем канал после того, как все горутины завершились

	fmt.Println("Scan completed.")
}

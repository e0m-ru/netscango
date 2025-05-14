package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const (
	MAX_PORT = 65535
)

var (
	timeout      = 333 * time.Millisecond
	workersCount = 1000
	portsRange   = []int{80, 443, 515, 631}
)

func main() {
	firstIP, lastIP := choicIPrange()
	targetsChan := make(chan string)
	resultsChan := make(chan string)
	var wg sync.WaitGroup
	for i := range workersCount {
		wg.Add(1)
		go worker(i, &wg, targetsChan, resultsChan)
	}
	go func(ip, lastIP net.IP) {
		for ; !ip.Equal(lastIP); incIP(ip) {
			for _, port := range portsRange {
				targetsChan <- fmt.Sprintf("%s:%d", firstIP.String(), port)
			}
		}
		close(targetsChan)
	}(firstIP, lastIP)
	go func() {
		for result := range resultsChan {
			fmt.Printf("%v\n", result)
		}
	}()
	wg.Wait()
	close(resultsChan)

	beep()
}

func choicIPrange() (net.IP, net.IP) {
	var networks = []*net.IPNet{}
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// Перебираем адреса интерфейса
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				networks = append(networks, ipNet)
			}
		}
	}
	for i, n := range networks {
		fmt.Printf("%d) %v\n", i+1, n)
	}
	fmt.Print("Select network number: ")
	var i int8
	fmt.Scanf("%d\n", &i)
	fmt.Printf("\rScan %v\n", networks[i-1])
	return getIPRange(networks[i-1])
}

func getIPRange(ipNet *net.IPNet) (net.IP, net.IP) {
	// Маскируем IP, чтобы получить начальный адрес сети
	firstIP := ipNet.IP.Mask(ipNet.Mask)

	// Создаём копию IP для вычисления последнего адреса
	lastIP := make(net.IP, len(firstIP))
	copy(lastIP, firstIP)

	// Применяем инвертированную маску к последнему IP
	mask := ipNet.Mask
	for i := 0; i < len(lastIP); i++ {
		lastIP[i] |= ^mask[i]
	}

	return firstIP, lastIP
}

func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

func beep() {
	fmt.Print("\a") // Пробуем стандартный BEL
	// Fallback для Windows
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "[console]::beep(800, 200)").Run()
	}
}

func worker(i int, wg *sync.WaitGroup, tasks, results chan string) {
	defer wg.Done()
	for task := range tasks {
		conn, err := net.DialTimeout("tcp", task, timeout)
		if err != nil {
			continue
		} else {
			results <- task
		}
		conn.Close()
	}
}

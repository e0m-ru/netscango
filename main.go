package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const (
	MAX_PORT = 65535
)

var (
	portsRange        []int
	workersCount      int
	timeout           int
	portString        string
	err               error
	defaultPortString = fmt.Sprintf("1-%d", MAX_PORT)
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		flag.PrintDefaults()
	}
	flag.IntVar(&workersCount, "w", 1000, "Determines the number of workers")
	flag.IntVar(&timeout, "t", 1000, "Determines the timeout for connection in milliseconds")
	flag.StringVar(&portString, "p", defaultPortString, "Ports define like -p [8080 || 1-1024 || 1-80,443,21-22,4455]")
	flag.Parse()
	// Parse port range
	portsRange, err = ParsePortRanges(portString)
	if err != nil {
		log.Fatalf("Error parsing ports: %v", err)
	}
	RUN()
}

func RUN() {
	targetsChan := make(chan string)
	resultsChan := make(chan string)
	firstIP, lastIP := choiceIPrange()
	var wg sync.WaitGroup
	for range workersCount {
		wg.Add(1)
		go worker(&wg, targetsChan, resultsChan)
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

func choiceIPrange() (net.IP, net.IP) {
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

func worker(wg *sync.WaitGroup, tasks, results chan string) {
	defer wg.Done()
	for task := range tasks {
		conn, err := net.DialTimeout("tcp", task, time.Millisecond*time.Duration(timeout))
		if err != nil {
			continue
		} else {
			results <- task
		}
		conn.Close()
	}
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Print("\033[H\033[2J")
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Очистка экрана не поддерживается")
	}
}

func beep() {
	fmt.Print("\a") // Пробуем стандартный BEL
	// Fallback для Windows
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "[console]::beep(800, 200)").Run()
	}
}

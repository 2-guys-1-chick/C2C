package client

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/2-guys-1-chick/c2c/cfg"
)

func RoundupConnect(loadedIp []net.IP) error {
	ipv4s, err := getOwnIPv4s()
	if err != nil {
		return err
	}

	tryoutIps := make(chan net.IP, 255*255)

	//var validIps []net.IP
	const workerCount = 100
	var wg sync.WaitGroup
	wg.Add(workerCount)

	timeout := time.Second
	for i := 0; i < workerCount; i++ { // 100 to constant
		go func() {
			for {
				tryoutIp, isOpen := <-tryoutIps
				if !isOpen {
					wg.Done()
					break
				}

				tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", tryoutIp.String(), cfg.GetPort()))
				if err != nil {
					continue
				}
				conn, err := net.DialTimeout("tcp", tcpAddr.String(), timeout)
				if err != nil {
					continue
				}

				fmt.Println("found new connection")
				handleNewConnection(conn)
			}
		}()
	}

	for _, ip := range ipv4s {
		for i := 0; i < 256; i++ {
			tryoutIp := make(net.IP, len(ip))
			copy(tryoutIp, ip)
			tryoutIp[3] = byte(i)

			if isInSlice(tryoutIp, ipv4s) {
				continue
			}

			if isInSlice(tryoutIp, loadedIp) {
				continue
			}

			tryoutIps <- tryoutIp
		}
	}

	close(tryoutIps)

	wg.Wait()

	return nil
}

func getOwnIPv4s() ([]net.IP, error) {
	ips, err := getOwnIPs()
	if err != nil {
		return nil, err
	}

	// easier to create new slice
	var ipv4s []net.IP
	for _, ip := range ips {
		if ip.IsLoopback() {
			continue
		}

		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}

		ipv4s = append(ipv4s, ipv4)
	}

	return ipv4s, err
}

func getOwnIPs() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			// handle error
			log.Println(err)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if !ip.IsUnspecified() {
				ips = append(ips, ip)
			}
		}
	}

	return ips, nil
}

func isInSlice(ip net.IP, ips []net.IP) bool {
	for _, i := range ips {
		if ip.Equal(i) {
			return true
		}
	}

	return false
}

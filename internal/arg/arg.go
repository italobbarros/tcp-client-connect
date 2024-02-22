// Pacote arg
package arg

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Config é uma estrutura para armazenar os valores dos argumentos
type Config struct {
	Addr []string
	Help bool
	Num  int
	Mode string
}

// ParseFlags processa os argumentos da linha de comando e retorna uma Config
func ParseFlags() (Config, error) {
	addrPtr := flag.String("addr", "", "address to be connected - <ip>:<port>")
	helpPtr := flag.Bool("help", false, "List all comands")
	numPtr := flag.Int("num", 1, "Number of connections")
	//modePtr := flag.String("mode", "", "Function Mode - Client or Server (not implemented)")

	// Parse os flags
	flag.Parse()

	config := Config{
		Help: *helpPtr,
		Num:  *numPtr,
		//Mode: *modePtr,
	}
	if config.Help {
		fmt.Println("Usage: tcpclient [options]")
		flag.PrintDefaults()
		return config, fmt.Errorf("")
	}

	if err := checkAddr(addrPtr, &config, *numPtr); err != nil {
		return config, err
	}

	return config, nil
}

func checkAddr(addr *string, cfg *Config, number int) error {
	addrList := strings.Split(*addr, ",")
	if len(addrList) == 1 {
		if err := isValidAddr(addrList[0]); err != nil {
			return err
		}
		var i int
		for i = 0; i < number; i++ {
			cfg.Addr = append(cfg.Addr, addrList[0])
		}
		return nil
	}
	var ip string
	for i, addr := range addrList {
		if i == 0 {
			if err := isValidAddr(addr); err != nil {
				return err
			}
			ip = strings.Split(addr, ":")[0]
			for i = 0; i < number; i++ {
				cfg.Addr = append(cfg.Addr, addr)
			}
			continue
		}
		addrParts := strings.Split(addr, ":")
		if len(addrParts) != 2 {
			return fmt.Errorf("invalid address. The format should be <IP>:<Port> \n Ex: localhost:8080")
		}
		if addrParts[0] == "" {
			if err := isValidPort(addrParts[1]); err != nil {
				return fmt.Errorf("addr pos: %d - %s", i, err.Error())
			}
			for i = 0; i < number; i++ {
				cfg.Addr = append(cfg.Addr, ip+addr)
			}
		} else {
			if err := isValidAddr(addr); err != nil {
				return err
			}
			ip = strings.Split(addr, ":")[0]
			for i = 0; i < number; i++ {
				cfg.Addr = append(cfg.Addr, addr)
			}
		}
	}
	return nil
}

func isValidAddr(addr string) error {
	if addr == "" {
		return fmt.Errorf("--addr <ip>:<port> not found \n Ex: --addr localhost:8080")
	}
	addrParts := strings.Split(addr, ":")
	if len(addrParts) != 2 {
		return fmt.Errorf("invalid address. The format should be <IP>:<Port> \n Ex: localhost:8080")
	}

	ip := addrParts[0]
	port := addrParts[1]

	if ip != "localhost" && net.ParseIP(ip) == nil {
		return fmt.Errorf("Invalid IP. The IP address should be IPv4 form.\n Ex: 10.10.10.180")
	}
	if err := isValidPort(port); err != nil {
		return err
	}
	return nil
}
func isValidPort(port string) error {
	v, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("Invalid port. The port should be a integer number between 1 and 65535. \n Ex: 8080")
	}
	if v < 1 || v > 65535 {
		return fmt.Errorf("%d is invalid port. The port should be a integer number between 1 and 65535. \n Ex: 8080", v)
	}
	return nil
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/app/hello", func(c fiber.Ctx) error {
		// 获取服务器IP地址
		info, err := getServerInfo()
		if err != nil {
			return c.SendString(err.Error())
		}

		resp := `
Hello world
Server IP Address: %s
Hostname: %s
`
		return c.SendString(fmt.Sprintf(resp, info.ip, info.hostname))
	})

	app.Get("/app/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(":8080"))
}

type serverInfo struct {
	hostname string
	ip       string
}

// 获取服务器信息
func getServerInfo() (*serverInfo, error) {
	// 获取本地主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// 获取主机地址列表
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}

	// 返回第一个非环回地址
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && !ip.IsLoopback() {
			return &serverInfo{
				hostname: hostname,
				ip:       addr,
			}, nil
		}
	}

	return nil, fmt.Errorf("No non-loopback address found")
}

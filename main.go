package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/luthermonson/go-proxmox"
)

var (
	domain   string
	apiLogin string
	apiKey   string
	pveName  string
)

func loadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	domain = os.Getenv("DOMAIN")
	apiLogin = os.Getenv("APILOGIN")
	apiKey = os.Getenv("APIKEY")
	pveName = os.Getenv("PVENAME")
}

func pveInfo(client *proxmox.Client, pveName string) {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		panic(err)
	}
	fmt.Println(pve.Name)
	fmt.Println(pve.Uptime)
	fmt.Println(pve.CPUInfo.CPUs)
	fmt.Println(pve.CPUInfo.Model)
}

func main() {
	loadENV()
	URL := &url.URL{}
	URL.Scheme = "https"
	URL.Host = domain
	URL.Path = "/api2/json"

	fmt.Println(URL)

	client := proxmox.NewClient(URL.String(),
		proxmox.WithAPIToken(apiLogin, apiKey),
	)

	pveInfo(client, pveName)
}

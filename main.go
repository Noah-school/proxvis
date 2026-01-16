package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/luthermonson/go-proxmox"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	domain := os.Getenv("DOMAIN")
	apiLogin := os.Getenv("APILOGIN")
	apiKey := os.Getenv("APIKEY")

	URL := "https://" + domain + "/api2/json"

	client := proxmox.NewClient(URL,
		proxmox.WithAPIToken(apiLogin, apiKey),
	)

	version, err := client.Version(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(version.Release)
}

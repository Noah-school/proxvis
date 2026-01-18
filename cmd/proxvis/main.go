package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/luthermonson/go-proxmox"
	"github.com/noah-school/proxvis/internal/config"
	"github.com/noah-school/proxvis/internal/pve"
)

func main() {
	cfg, err := config.LoadENV()
	if err != nil {
		log.Fatal(err)
	}

	URL := &url.URL{
		Scheme: "https",
		Host:   cfg.Domain,
		Path:   "/api2/json",
	}

	fmt.Println(URL)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.TLSInsecure,
			},
		},
	}

	client := proxmox.NewClient(URL.String(),
		proxmox.WithAPIToken(cfg.APILogin, cfg.APIKey),
		proxmox.WithHTTPClient(httpClient),
	)

	manager := &pve.Manager{Client: client}
	ctx := context.Background()

	if err := manager.PveInfo(ctx, cfg.PVEName); err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------------------------------")

	fmt.Println("Listing ISOs...")
	isos, err := manager.ListISOs(ctx, cfg.PVEName)
	if err != nil {
		log.Println("Error listing ISOs:", err)
	}
	var selectedISO string
	if len(isos) > 0 {
		fmt.Println("Found ISOs:")
		for _, iso := range isos {
			fmt.Println("- " + iso)
		}
		selectedISO = isos[1]
		fmt.Println("Selected ISO:", selectedISO)
	} else {
		fmt.Println("No ISOs found.")
	}

	fmt.Println("make 404 vm")
	if err := manager.MakeVM(ctx, cfg.PVEName, 404, pve.VMConfig{
		Name:    "test-vm-404",
		Memory:  2048,
		Cores:   2,
		Sockets: 1,
		OnBoot:  true,
		ISO:     selectedISO,
	}); err != nil {
		log.Fatal(err)
	}
	if err := manager.StartVM(ctx, cfg.PVEName, 404); err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting 404 vm")

	if err := manager.WaitForStatus(ctx, cfg.PVEName, 404, "running"); err != nil {
		log.Fatal(err)
	}
	if err := manager.ShowVM(ctx, cfg.PVEName, 404); err != nil {
		log.Println(err)
	}

	fmt.Println("stopping 404 vm")
	if err := manager.StopVM(ctx, cfg.PVEName, 404); err != nil {
		log.Fatal(err)
	}

	if err := manager.WaitForStatus(ctx, cfg.PVEName, 404, "stopped"); err != nil {
		log.Fatal(err)
	}
	if err := manager.ShowVM(ctx, cfg.PVEName, 404); err != nil {
		log.Println(err)
	}

	fmt.Println("delete 404 vm")
	if err := manager.DelVM(ctx, cfg.PVEName, 404); err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------------------------------")
	fmt.Println("Generating Topology:")
	topo, err := manager.GenerateTopology(ctx, cfg.PVEName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(topo)
}

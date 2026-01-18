package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/luthermonson/go-proxmox"
)

type App struct {
	Client *proxmox.Client
}

type Config struct {
	Domain      string
	APILogin    string
	APIKey      string
	PVEName     string
	TLSInsecure bool
}

func loadENV() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	tlsInsecure, _ := strconv.ParseBool(os.Getenv("TLS_INSECURE"))

	return &Config{
		Domain:      os.Getenv("DOMAIN"),
		APILogin:    os.Getenv("APILOGIN"),
		APIKey:      os.Getenv("APIKEY"),
		PVEName:     os.Getenv("PVENAME"),
		TLSInsecure: tlsInsecure,
	}, nil
}

func (app *App) selectNode(ctx context.Context, pveName string) (*proxmox.Node, error) {
	pve, err := app.Client.Node(ctx, pveName)
	if err != nil {
		return nil, err
	}
	return pve, nil
}

func (app *App) pveInfo(ctx context.Context, pveName string) error {
	pve, err := app.Client.Node(ctx, pveName)
	if err != nil {
		return err
	}
	fmt.Println(pve.Name)
	fmt.Println(pve.Uptime)
	fmt.Println(pve.CPUInfo.CPUs)
	fmt.Println(pve.CPUInfo.Model)
	return nil
}

func (app *App) makeVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := app.Client.Node(ctx, pveName)
	if err != nil {
		return err
	}
	newVM, err := pve.NewVirtualMachine(ctx, VMID)
	if err != nil {
		return err
	}
	fmt.Println(newVM.ID)
	return nil
}

func (app *App) delVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := app.selectNode(ctx, pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(ctx, VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Stop(ctx); err != nil {
		return err
	}
	if _, err := vm.Delete(ctx); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func (app *App) showVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := app.selectNode(ctx, pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(ctx, VMID)
	if err != nil {
		return err
	}
	fmt.Println(vm.Status)
	return nil
}

func (app *App) startVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := app.selectNode(ctx, pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(ctx, VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (app *App) stopVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := app.selectNode(ctx, pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(ctx, VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (app *App) waitForStatus(ctx context.Context, pveName string, VMID int, status string) error {
	fmt.Printf("Waiting for VM %d to be %s...\n", VMID, status)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			pve, err := app.selectNode(ctx, pveName)
			if err != nil {
				return err
			}
			vm, err := pve.VirtualMachine(ctx, VMID)
			if err != nil {
				return err
			}
			if vm.Status == status {
				return nil
			}
		}
	}
}

func main() {
	config, err := loadENV()
	if err != nil {
		log.Fatal(err)
	}

	URL := &url.URL{
		Scheme: "https",
		Host:   config.Domain,
		Path:   "/api2/json",
	}

	fmt.Println(URL)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.TLSInsecure,
			},
		},
	}

	client := proxmox.NewClient(URL.String(),
		proxmox.WithAPIToken(config.APILogin, config.APIKey),
		proxmox.WithHTTPClient(httpClient),
	)

	app := &App{Client: client}
	ctx := context.Background()

	if err := app.pveInfo(ctx, config.PVEName); err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------------------------------")
	fmt.Println("make 404 vm")
	if err := app.makeVM(ctx, config.PVEName, 404); err != nil {
		log.Fatal(err)
	}
	if err := app.startVM(ctx, config.PVEName, 404); err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting 404 vm")

	if err := app.waitForStatus(ctx, config.PVEName, 404, "running"); err != nil {
		log.Fatal(err)
	}
	if err := app.showVM(ctx, config.PVEName, 404); err != nil {
		log.Println(err)
	}

	fmt.Println("stopping 404 vm")
	if err := app.stopVM(ctx, config.PVEName, 404); err != nil {
		log.Fatal(err)
	}

	if err := app.waitForStatus(ctx, config.PVEName, 404, "stopped"); err != nil {
		log.Fatal(err)
	}
	if err := app.showVM(ctx, config.PVEName, 404); err != nil {
		log.Println(err)
	}

	fmt.Println("delete 404 vm")
	if err := app.delVM(ctx, config.PVEName, 404); err != nil {
		log.Fatal(err)
	}
}

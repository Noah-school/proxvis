package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/luthermonson/go-proxmox"
)

type App struct {
	Client *proxmox.Client
}

func loadENV() (string, string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	domain := os.Getenv("DOMAIN")
	apiLogin := os.Getenv("APILOGIN")
	apiKey := os.Getenv("APIKEY")
	pveName := os.Getenv("PVENAME")

	return domain, apiLogin, apiKey, pveName
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

func main() {
	domain, apiLogin, apiKey, pveName := loadENV()
	URL := &url.URL{}
	URL.Scheme = "https"
	URL.Host = domain
	URL.Path = "/api2/json"

	fmt.Println(URL)

	client := proxmox.NewClient(URL.String(),
		proxmox.WithAPIToken(apiLogin, apiKey),
	)

	app := &App{Client: client}
	ctx := context.Background()

	if err := app.pveInfo(ctx, pveName); err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------------------------------")
	fmt.Println("make 404 vm")
	if err := app.makeVM(ctx, pveName, 404); err != nil {
		log.Fatal(err)
	}
	if err := app.startVM(ctx, pveName, 404); err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting 404 vm")
	time.Sleep(4 * time.Second)
	if err := app.showVM(ctx, pveName, 404); err != nil {
		log.Println(err)
	}
	fmt.Println("stopping 404 vm")
	if err := app.stopVM(ctx, pveName, 404); err != nil {
		log.Fatal(err)
	}
	time.Sleep(4 * time.Second)
	if err := app.showVM(ctx, pveName, 404); err != nil {
		log.Println(err)
	}
	fmt.Println("delete 404 vm")
	time.Sleep(4 * time.Second)
	if err := app.delVM(ctx, pveName, 404); err != nil {
		log.Fatal(err)
	}
}

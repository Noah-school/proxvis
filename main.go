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

var (
	client *proxmox.Client
)

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

func selectNode(pveName string) (*proxmox.Node, error) {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		return nil, err
	}
	return pve, nil
}

func pveInfo(pveName string) error {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		return err
	}
	fmt.Println(pve.Name)
	fmt.Println(pve.Uptime)
	fmt.Println(pve.CPUInfo.CPUs)
	fmt.Println(pve.CPUInfo.Model)
	return nil
}

func makeVM(pveName string, VMID int) error {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		return err
	}
	newVM, err := pve.NewVirtualMachine(context.Background(), VMID)
	if err != nil {
		return err
	}
	fmt.Println(newVM.ID)
	return nil
}

func delVM(pveName string, VMID int) error {
	pve, err := selectNode(pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Stop(context.Background()); err != nil {
		return err
	}
	if _, err := vm.Delete(context.Background()); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func showVM(pveName string, VMID int) error {
	pve, err := selectNode(pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		return err
	}
	fmt.Println(vm.Status)
	return nil
}

func startVM(pveName string, VMID int) error {
	pve, err := selectNode(pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Start(context.Background()); err != nil {
		return err
	}
	return nil
}

func stopVM(pveName string, VMID int) error {
	pve, err := selectNode(pveName)
	if err != nil {
		return err
	}
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		return err
	}
	if _, err := vm.Stop(context.Background()); err != nil {
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

	client = proxmox.NewClient(URL.String(),
		proxmox.WithAPIToken(apiLogin, apiKey),
	)

	if err := pveInfo(pveName); err != nil {
		log.Fatal(err)
	}
	fmt.Println("----------------------------------")
	fmt.Println("make 404 vm")
	if err := makeVM(pveName, 404); err != nil {
		log.Fatal(err)
	}
	if err := startVM(pveName, 404); err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting 404 vm")
	time.Sleep(4 * time.Second)
	if err := showVM(pveName, 404); err != nil {
		log.Println(err)
	}
	fmt.Println("stopping 404 vm")
	if err := stopVM(pveName, 404); err != nil {
		log.Fatal(err)
	}
	time.Sleep(4 * time.Second)
	if err := showVM(pveName, 404); err != nil {
		log.Println(err)
	}
	fmt.Println("delete 404 vm")
	time.Sleep(4 * time.Second)
	if err := delVM(pveName, 404); err != nil {
		log.Fatal(err)
	}
}

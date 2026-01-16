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

func selectNode(client *proxmox.Client) *proxmox.Node {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		panic(err)
	}
	return pve
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

func makeVM(client *proxmox.Client, VMID int) {
	pve, err := client.Node(context.Background(), pveName)
	if err != nil {
		panic(err)
	}
	newVM, err := pve.NewVirtualMachine(context.Background(), VMID)
	if err != nil {
		panic(err)
	}
	fmt.Println(newVM.ID)
}

func delVM(client *proxmox.Client, VMID int) {
	pve := selectNode(client)
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		panic(err)
	}
	vm.Stop(context.Background())
	vm.Delete(context.Background())
	fmt.Println()
}

func showVM(client *proxmox.Client, VMID int) {
	pve := selectNode(client)
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		panic(err)
	}
	fmt.Println(vm.Status)
}

func startVM(client *proxmox.Client, VMID int) {
	pve := selectNode(client)
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		panic(err)
	}
	vm.Start(context.Background())
}

func stopVM(client *proxmox.Client, VMID int) {
	pve := selectNode(client)
	vm, err := pve.VirtualMachine(context.Background(), VMID)
	if err != nil {
		panic(err)
	}
	vm.Stop(context.Background())
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

	fmt.Println("make 404 vm")
	makeVM(client, 404)
	startVM(client, 404)
	fmt.Println("starting 404 vm")
	time.Sleep(4 * time.Second)
	showVM(client, 404)
	fmt.Println("stopping 404 vm")
	stopVM(client, 404)
	time.Sleep(4 * time.Second)
	showVM(client, 404)
	fmt.Println("delete 404 vm")
	time.Sleep(4 * time.Second)
	delVM(client, 404)
}

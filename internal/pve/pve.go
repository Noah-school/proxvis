package pve

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/luthermonson/go-proxmox"
)

type Manager struct {
	Client *proxmox.Client
}

type VMConfig struct {
	Name    string
	Memory  int
	Cores   int
	Sockets int
	OnBoot  bool
	ISO     string
}

func (m *Manager) selectNode(ctx context.Context, pveName string) (*proxmox.Node, error) {
	pve, err := m.Client.Node(ctx, pveName)
	if err != nil {
		return nil, err
	}
	return pve, nil
}

func (m *Manager) PveInfo(ctx context.Context, pveName string) error {
	pve, err := m.Client.Node(ctx, pveName)
	if err != nil {
		return err
	}
	fmt.Println(pve.Name)
	fmt.Println(pve.Uptime)
	fmt.Println(pve.CPUInfo.CPUs)
	fmt.Println(pve.CPUInfo.Model)
	return nil
}

func (m *Manager) ListISOs(ctx context.Context, pveName string) ([]string, error) {
	pve, err := m.Client.Node(ctx, pveName)
	if err != nil {
		return nil, err
	}

	storages, err := pve.Storages(ctx)
	if err != nil {
		return nil, err
	}

	var isos []string
	for _, storage := range storages {
		if strings.Contains(storage.Content, "iso") {
			content, err := storage.GetContent(ctx)
			if err != nil {
				continue
			}
			for _, item := range content {
				if strings.Contains(item.Volid, "iso/") {
					isos = append(isos, item.Volid)
				}
			}
		}
	}
	return isos, nil
}

func (m *Manager) MakeVM(ctx context.Context, pveName string, VMID int, config VMConfig) error {
	pve, err := m.Client.Node(ctx, pveName)
	if err != nil {
		return err
	}

	var options []proxmox.VirtualMachineOption
	if config.Name != "" {
		options = append(options, proxmox.VirtualMachineOption{Name: "name", Value: config.Name})
	}
	if config.Memory > 0 {
		options = append(options, proxmox.VirtualMachineOption{Name: "memory", Value: config.Memory})
	}
	if config.Cores > 0 {
		options = append(options, proxmox.VirtualMachineOption{Name: "cores", Value: config.Cores})
	}
	if config.Sockets > 0 {
		options = append(options, proxmox.VirtualMachineOption{Name: "sockets", Value: config.Sockets})
	}
	if config.OnBoot {
		options = append(options, proxmox.VirtualMachineOption{Name: "onboot", Value: 1})
	}
	if config.ISO != "" {
		options = append(options, proxmox.VirtualMachineOption{Name: "ide2", Value: config.ISO + ",media=cdrom"})
	}

	task, err := pve.NewVirtualMachine(ctx, VMID, options...)
	if err != nil {
		return err
	}

	fmt.Printf("VM Creation Task: %v\n", task.UPID)

	return nil
}

func (m *Manager) DelVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := m.selectNode(ctx, pveName)
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

func (m *Manager) ShowVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := m.selectNode(ctx, pveName)
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

func (m *Manager) StartVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := m.selectNode(ctx, pveName)
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

func (m *Manager) StopVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := m.selectNode(ctx, pveName)
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

func (m *Manager) WaitForStatus(ctx context.Context, pveName string, VMID int, status string) error {
	fmt.Printf("Waiting for VM %d to be %s...\n", VMID, status)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			pve, err := m.selectNode(ctx, pveName)
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

package pve

import (
	"context"
	"fmt"
	"time"

	"github.com/luthermonson/go-proxmox"
)

type Manager struct {
	Client *proxmox.Client
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

func (m *Manager) MakeVM(ctx context.Context, pveName string, VMID int) error {
	pve, err := m.Client.Node(ctx, pveName)
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

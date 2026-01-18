package pve

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func (m *Manager) GenerateTopology(ctx context.Context, pveName string) (string, error) {
	node, err := m.selectNode(ctx, pveName)
	if err != nil {
		return "", err
	}

	networks, err := node.Networks(ctx)
	if err != nil {
		return "", err
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Iface < networks[j].Iface
	})

	var sb strings.Builder
	sb.WriteString("graph TD\n")

	sb.WriteString("    classDef bridge fill:#fff9c4,stroke:#fbc02d,stroke-width:2px;\n")
	sb.WriteString("    classDef iface fill:#f3e5f5,stroke:#4a148c,stroke-width:2px;\n")
	sb.WriteString("    classDef vm fill:#e1f5fe,stroke:#0277bd,stroke-width:2px;\n")

	sb.WriteString(fmt.Sprintf("    subgraph %s [\"Node: %s\"]\n", cleanID(pveName), pveName))
	sb.WriteString("        direction TB\n")
	sb.WriteString("        style " + cleanID(pveName) + " fill:#ffffff,stroke:#333,stroke-width:2px\n")

	var bridgeNodes, ifaceNodes, vmNodes strings.Builder
	var links strings.Builder

	knownBridges := make(map[string]bool)

	for _, net := range networks {
		netID := cleanID(net.Iface)
		label := fmt.Sprintf("<b>%s</b><br/>%s", net.Iface, net.Type)
		if net.CIDR != "" {
			label += fmt.Sprintf("<br/>%s", net.CIDR)
		}

		if net.Type == "bridge" {
			knownBridges[net.Iface] = true
			bridgeNodes.WriteString(fmt.Sprintf("            %s[\"%s\"]:::bridge\n", netID, label))

			if net.BridgePorts != "" {
				rawPorts := strings.FieldsFunc(net.BridgePorts, func(r rune) bool {
					return r == ' ' || r == ','
				})
				for _, port := range rawPorts {
					port = strings.TrimSpace(port)
					if port != "" {
						portID := cleanID(port)
						links.WriteString(fmt.Sprintf("    %s --> %s\n", portID, netID))
					}
				}
			}
		} else {
			ifaceNodes.WriteString(fmt.Sprintf("            %s[\"%s\"]:::iface\n", netID, label))
		}
	}

	vms, err := node.VirtualMachines(ctx)
	if err != nil {
		return "", err
	}

	sort.Slice(vms, func(i, j int) bool {
		idI, _ := strconvToInt(vms[i].VMID)
		idJ, _ := strconvToInt(vms[j].VMID)
		return idI < idJ
	})

	for _, vmSummary := range vms {
		vmIDInt, _ := strconvToInt(vmSummary.VMID)
		vm, err := node.VirtualMachine(ctx, vmIDInt)
		if err != nil {
			fmt.Printf("Warning: failed to fetch details for VM %v\n", vmSummary.VMID)
			continue
		}

		if vm.VirtualMachineConfig == nil {
			continue
		}

		vmNodeID := fmt.Sprintf("vm_%v", vm.VMID)
		vmLabel := fmt.Sprintf("VM %v<br/>%s", vm.VMID, vm.Name)
		vmNodes.WriteString(fmt.Sprintf("            %s[\"%s\"]:::vm\n", vmNodeID, vmLabel))

		nets := vm.VirtualMachineConfig.MergeNets()
		var netKeys []string
		for k := range nets {
			netKeys = append(netKeys, k)
		}
		sort.Strings(netKeys)

		for _, k := range netKeys {
			netConf := nets[k]
			bridge := parseBridge(netConf)
			if bridge != "" && knownBridges[bridge] {
				bridgeID := cleanID(bridge)
				links.WriteString(fmt.Sprintf("    %s -.-> %s\n", vmNodeID, bridgeID))
			}
		}
	}

	sb.WriteString("        subgraph Cluster_VMs [\"Virtual Machines\"]\n")
	sb.WriteString("            direction LR\n")
	sb.WriteString("            style Cluster_VMs fill:none,stroke:none\n")
	sb.WriteString(vmNodes.String())
	sb.WriteString("        end\n")

	sb.WriteString("        subgraph Cluster_Net [\"Network\"]\n")
	sb.WriteString("            style Cluster_Net fill:none,stroke:none\n")

	sb.WriteString("            subgraph Cluster_Bridges [\"Bridges\"]\n")
	sb.WriteString("                style Cluster_Bridges fill:none,stroke:none\n")
	sb.WriteString(bridgeNodes.String())
	sb.WriteString("            end\n")

	sb.WriteString("            subgraph Cluster_Ifaces [\"Interfaces\"]\n")
	sb.WriteString("                style Cluster_Ifaces fill:none,stroke:none\n")
	sb.WriteString(ifaceNodes.String())
	sb.WriteString("            end\n")

	sb.WriteString("        end\n")

	sb.WriteString("    end\n")

	sb.WriteString(links.String())

	return sb.String(), nil
}

func cleanID(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func parseBridge(conf string) string {
	parts := strings.Split(conf, ",")
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "bridge=") {
			return strings.TrimPrefix(strings.TrimSpace(part), "bridge=")
		}
	}
	return ""
}

func strconvToInt(v interface{}) (int, error) {
	s := fmt.Sprint(v)
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

# Proxvis
Derived from Proxmox Visualizer, Proxvis aims to allow configuration of your hypervisor with a visualized topology dashboard using the API.

## Features
- **Visual Topology**: Generate a visual representation of your Proxmox infrastructure.
- **VM Management**: Create, start, stop, and delete VMs programmatically.

## Prerequisites
- **Go**: Version 1.24 or later.
- **Proxmox VE**: A running Proxmox VE instance.
- **Proxmox API Token**: An API token with appropriate permissions.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/noah-school/proxvis.git
   cd proxvis
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

1. Create a `.env` file from the example:
   ```bash
   cp .env_example .env
   ```

2. Edit `.env` with your Proxmox credentials:
   ```env
   DOMAIN=pve.domain.com                            # Proxmox host
   APILOGIN=user@pam!ID                             # Token ID
   APIKEY=00000000-0000-0000-0000-000000000000      # Your API Secret Key
   PVENAME=mini                                     # The name of your Proxmox node
   ```

   > **Note**: Ensure the API Token has permissions to read node status, storage, and manage VMs.

## Usage

To run the application:

```bash
go run cmd/proxvis/main.go
```

This will currently:
1. Connect to your Proxmox instance.
2. List available ISO images.
3. Perform specific VM operations (Create -> Start -> Stop -> Delete) as a demonstration.
4. Generate a topology of your Proxmox node.

## Development

Project structure:
- `cmd/proxvis`: Main application.
- `internal/config`: Configuration management.
- `internal/pve`: Proxmox VE interaction logic.
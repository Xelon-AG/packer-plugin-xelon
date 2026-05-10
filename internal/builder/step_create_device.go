package builder

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"

	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var _ multistep.Step = (*stepCreateDevice)(nil)

// stepCreateDevice represents a build step that creates Xelon device.
type stepCreateDevice struct {
	deviceID string
}

func (s *stepCreateDevice) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*xelon.Client)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating a Xelon device...")

	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	createRequest := &xelon.DeviceCreateRequest{
		CPUCores:             2,
		RAM:                  2,
		DiskSize:             10,
		DisplayName:          name,
		HostName:             name,
		Password:             config.RootPassword,
		PasswordConfirmation: config.RootPassword,
		SwapDiskSize:         1,
		TemplateID:           config.SourceTemplateID,
		TenantID:             config.TenantID,
		Networks: []xelon.DeviceCreateNetwork{
			{
				ConnectOnPowerOn: true,
				NetworkID:        config.NetworkID,
			},
		},
	}
	sshKeyID := state.Get("ssh_key_id").(string)
	createRequest.SSHKeyID = sshKeyID

	log.Printf("[DEBUG] Creating Xelon device: %v", createRequest)
	device, _, err := client.Devices.Create(ctx, createRequest)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Sayf("Device ID: %s", device.ID)
	s.deviceID = device.ID

	ui.Sayf("Waiting for device (%s) to become ready...", device.ID)
	err = waitDevicePowerStateOn(ctx, client, device.ID)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	err = waitDeviceStateReady(ctx, client, device.ID)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	deviceIPAddress := ""
	deviceNetworks, _, err := client.Devices.GetNetworkInfo(ctx, device.ID)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	for _, deviceNetwork := range deviceNetworks {
		if deviceNetwork.Connected && deviceNetwork.ID == config.NetworkID {
			if len(deviceNetwork.IPAddresses) == 0 {
				continue
			}

			for _, ipAddress := range deviceNetwork.IPAddresses {
				if ipAddress.Is4() {
					deviceIPAddress = ipAddress.String()
					break
				}
			}
		}
	}

	state.Put("device_id", device.ID)
	state.Put("device_ip", deviceIPAddress)

	return multistep.ActionContinue
}

func (s *stepCreateDevice) Cleanup(state multistep.StateBag) {
	client := state.Get("client").(*xelon.Client)
	_ = state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if s.deviceID == "" {
		return
	}

	ctx := context.Background()
	device, resp, err := client.Devices.Get(ctx, s.deviceID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return
		}
		ui.Errorf("Error getting Xelon device (%s): %v", s.deviceID, err)
		return
	}

	if device.PoweredOn {
		_, err := client.Devices.Stop(ctx, s.deviceID)
		if err != nil {
			ui.Errorf("Error stopping Xelon device (%s): %v", s.deviceID, err)
			return
		}

		err = waitDevicePowerStateOff(ctx, client, s.deviceID)
		if err != nil {
			ui.Errorf("Error waiting for Xelon device (%s) to be powered off: %v", s.deviceID, err)
			return
		}
	}

	_, err = client.Devices.Delete(ctx, s.deviceID)
	if err != nil {
		ui.Errorf(" Error deleting Xelon device (%s): %v", s.deviceID, err)
	}
}

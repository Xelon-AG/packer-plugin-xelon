package builder

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	"golang.org/x/crypto/ssh"

	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var _ multistep.Step = (*stepCreateSSHKey)(nil)

// stepCreateSSHKey represents a build step that generates SSH key pairs.
type stepCreateSSHKey struct {
	Debug        bool
	DebugKeyPath string
	SSHKeyID     string

	doCleanup bool
}

func (s *stepCreateSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*xelon.Client)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if config.Comm.SSHPrivateKeyFile != "" {
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := config.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		config.Comm.SSHPrivateKey = privateKeyBytes
		return multistep.ActionContinue
	}

	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	ui.Sayf("Creating temporary SSH key: %s", name)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	config.Comm.SSHPrivateKey = pem.EncodeToMemory(&privateKeyBlock)
	config.Comm.SSHPublicKey = ssh.MarshalAuthorizedKey(publicKey)

	sshKey, _, err := client.SSHKeys.Create(ctx, &xelon.SSHKeyCreateRequest{
		SSHKey: xelon.SSHKey{
			Name:      name,
			PublicKey: string(config.Comm.SSHPublicKey),
		},
	})
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	s.doCleanup = true
	s.SSHKeyID = sshKey.ID
	state.Put("ssh_key_id", sshKey.ID)

	if s.Debug {
		ui.Sayf("Saving key for debug purposes: %s", s.DebugKeyPath)

		f, err := os.Create(s.DebugKeyPath)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		defer func() { _ = f.Close() }()

		err = pem.Encode(f, &privateKeyBlock)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if runtime.GOOS != "windows" {
			if err := f.Chmod(0600); err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}
		}
	}

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	client := state.Get("client").(*xelon.Client)
	ui := state.Get("ui").(packer.Ui)

	if !s.doCleanup {
		return
	}

	ui.Say("Deleting temporary SSH key...")
	_, err := client.SSHKeys.Delete(context.Background(), s.SSHKeyID)
	if err != nil {
		ui.Errorf("Error deleting temporary SSH key (%s). Please delete manually: %s", s.SSHKeyID, err)
	}

	if s.Debug {
		if err := os.Remove(s.DebugKeyPath); err != nil {
			ui.Errorf("Error removing debug temporary SSH key '%s': %s", s.DebugKeyPath, err)
		}
	}
}

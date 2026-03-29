package builder

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var (
	_ packer.Artifact = (*Artifact)(nil)
)

// Artifact is an artifact implementation that contains built Xelon Template.
type Artifact struct {
	Client *xelon.Client
	// StateData stores data such as GeneratedData to be shared with post-processors
	StateData map[string]any
	// TemplateID is the ID of Xelon Template
	TemplateID string
	// TemplateName is the name of Xelon Template
	TemplateName string
}

func (a *Artifact) BuilderId() string {
	return PluginBuilderID
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.TemplateID
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Xelon template: %s (%s)", a.TemplateName, a.TemplateID)
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	log.Printf("Deleting template: %s (%s)", a.TemplateName, a.TemplateID)
	_, err := a.Client.Templates.Delete(context.TODO(), a.TemplateID)
	return err
}

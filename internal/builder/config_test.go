package builder

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/stretchr/testify/assert"
)

func TestTemplateConfig_Prepare(t *testing.T) {
	type testCase struct {
		input               *TemplateConfig
		expectedName        string
		expectedDescription string
		interpolated        bool
	}
	tests := map[string]testCase{
		"default values": {
			input: &TemplateConfig{
				TemplateName: "",
			},
			interpolated:        true,
			expectedName:        "packer-",
			expectedDescription: "Created by Packer",
		},
		"custom values": {
			input: &TemplateConfig{
				TemplateName:        "template name",
				TemplateDescription: "template description",
			},
			interpolated:        false,
			expectedName:        "template name",
			expectedDescription: "template description",
		},
	}

	freezeTime(t, time.Unix(1234567890, 0))
	timestampStr := "1234567890"

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			errs := test.input.Prepare(interpolate.Context{}, nil)

			expectedName := test.expectedName
			if test.interpolated {
				expectedName = fmt.Sprintf("%s%s", test.expectedName, timestampStr)
			}

			assert.Nil(t, errs)
			assert.Equal(t, expectedName, test.input.TemplateName)
			assert.Equal(t, test.expectedDescription, test.input.TemplateDescription)
		})
	}
}

// sets {{timestamp}} Packer expression to pinned value `at`.
func freezeTime(t *testing.T, at time.Time) {
	t.Helper()
	orig := interpolate.InitTime
	interpolate.InitTime = at.UTC()
	t.Cleanup(func() { interpolate.InitTime = orig })
}

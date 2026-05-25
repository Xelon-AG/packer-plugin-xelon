package config

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/stretchr/testify/assert"
)

func TestAccessConfig_Prepare(t *testing.T) {
	type testCase struct {
		input          *AccessConfig
		expectedConfig AccessConfig
		expectedErrors int
	}
	tests := map[string]testCase{
		"required values": {
			input: &AccessConfig{
				ClientID: "client-id-value",
				Token:    "token-value",
			},
			expectedConfig: AccessConfig{
				ClientID: "client-id-value",
				Token:    "token-value",
			},
		},
		"missing client_id": {
			input: &AccessConfig{
				Token: "token-value",
			},
			expectedConfig: AccessConfig{},
			expectedErrors: 1,
		},
		"missing token": {
			input: &AccessConfig{
				ClientID: "client-id-value",
			},
			expectedConfig: AccessConfig{},
			expectedErrors: 1,
		},
		"missing client_id and token": {
			input:          &AccessConfig{},
			expectedConfig: AccessConfig{},
			expectedErrors: 2,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			errs := test.input.Prepare(&interpolate.Context{}, nil)
			actualConfig := test.expectedConfig

			if test.expectedErrors > 0 {
				assert.Len(t, errs.Errors, test.expectedErrors)
			} else {
				assert.Nil(t, errs)
				assert.Equal(t, test.expectedConfig, actualConfig)
			}
		})
	}
}

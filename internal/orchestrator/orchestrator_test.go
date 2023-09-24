package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrchestrator_isBlocked(t *testing.T) {
	var blockedDomains = []string{
		"localhost",
		"ngrok.io",
	}

	tests := []struct {
		domain string
		want   bool
	}{
		{
			domain: "localhost",
			want:   true,
		},
		{
			domain: "subdomain.localhost",
			want:   true,
		},
		{
			domain: "ngrok.io",
			want:   true,
		},
		{
			domain: "reallylong.toto.subdomain.ngrok.io",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			o := Orchestrator{
				config: OrchestratorConfig{
					BlockedDomains: blockedDomains,
				},
			}
			got := o.isBlocked(tt.domain)
			assert.Equal(t, tt.want, got)
		})
	}
}

package config

import (
	"testing"

	"kapigen.kateops.com/internal/pipeline/types"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name     string
		service  *Service
		expected error
	}{
		{
			name: "Valid service",
			service: &Service{
				Name:      "test-service",
				Port:      8080,
				ImageName: "test-image",
			},
			expected: nil,
		},
		{
			name: "Missing name",
			service: &Service{
				Port:      8080,
				ImageName: "test-image",
			},
			expected: types.NewMissingArgError("service.name"),
		},
		{
			name: "Invalid port",
			service: &Service{
				Name:      "test-service",
				Port:      0,
				ImageName: "test-image",
			},
			expected: types.DetailedErrorf("service: 'test-service', invalid port %d (must be 1 - 65535)", 0),
		},
		{
			name: "Missing image name and docker",
			service: &Service{
				Name: "test-service",
				Port: 8080,
			},
			expected: types.NewMissingArgsError("service.imageName", "service.docker"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.service.Validate()
			if err != nil && test.expected == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error: %v, but got nil", test.expected)
			} else if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}

func TestServices_Validate(t *testing.T) {
	tests := []struct {
		name     string
		services Services
		expected error
	}{
		{
			name: "Valid services",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "service2", Port: 8081, ImageName: "test-image"},
			},
			expected: nil,
		},
		{
			name: "Duplicate port",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "service2", Port: 8080, ImageName: "test-image"},
			},
			expected: types.DetailedErrorf("service: 'service2', referencing occupied port: %d", 8080),
		},
		{
			name: "Invalid service",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "", Port: 8081, ImageName: "test-image"},
			},
			expected: types.NewMissingArgError("service.name"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.services.Validate()
			if err != nil && test.expected == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error: %v, but got nil", test.expected)
			} else if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}

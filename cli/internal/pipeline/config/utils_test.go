package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name     string
		service  *Service
		expected string
	}{
		{
			name: "Valid service",
			service: &Service{
				Name:      "test-service",
				Port:      8080,
				ImageName: "test-image",
			},
			expected: "",
		},
		{
			name: "Missing name",
			service: &Service{
				Port:      8080,
				ImageName: "test-image",
			},
			expected: "service.name",
		},
		{
			name: "Invalid port",
			service: &Service{
				Name:      "test-service",
				Port:      0,
				ImageName: "test-image",
			},
			expected: fmt.Sprintf("service: 'test-service', invalid port %d (must be 1 - 65535)", 0),
		},
		{
			name: "Missing image name and docker",
			service: &Service{
				Name: "test-service",
				Port: 8080,
			},
			expected: "'service.imageName', 'service.docker'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.service.Validate()
			if err != nil && test.expected == "" {
				t.Errorf("Unexpected error: %v", err)
			} else if err == nil && test.expected != "" {
				t.Errorf("Expected error: %v, but got nil", test.expected)
			} else if err != nil && !strings.Contains(err.Error(), test.expected) {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}

func TestServices_Validate(t *testing.T) {
	tests := []struct {
		name     string
		services Services
		expected string
	}{
		{
			name: "Valid services",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "service2", Port: 8081, ImageName: "test-image"},
			},
			expected: "",
		},
		{
			name: "Duplicate port",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "service2", Port: 8080, ImageName: "test-image"},
			},
			expected: fmt.Sprintf("service: 'service2', referencing occupied port: %d", 8080),
		},
		{
			name: "Invalid service",
			services: Services{
				{Name: "service1", Port: 8080, ImageName: "test-image"},
				{Name: "", Port: 8081, ImageName: "test-image"},
			},
			expected: "service.name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.services.Validate()
			if err != nil && test.expected == "" {
				t.Errorf("Unexpected error: %v", err)
			} else if err == nil && test.expected != "" {
				t.Errorf("Expected error: %v, but got nil", test.expected)
			} else if err != nil && !strings.Contains(err.Error(), test.expected) {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}

func TestJobMode_String(t *testing.T) {
	t.Run("can find job mode permissive", func(t *testing.T) {
		expectation := "permissive"
		value, err := JobModeEnum.Value(Permissive)
		if err != nil {
			t.Error(err)
		}

		if value != expectation {
			t.Errorf("should be equal")
		}
		if JobModeEnum.ValueSafe(Permissive) != expectation {
			t.Errorf("should be equal")
		}
	})
	t.Run("can not find job mode", func(t *testing.T) {
		var mode JobMode = -1
		if JobModeEnum.ValueSafe(mode) != "" {
			t.Errorf("should be empty")
		}

		if _, err := JobModeEnum.Value(mode); err == nil {
			t.Error("should not be found return error, received nil")
		}
	})
}

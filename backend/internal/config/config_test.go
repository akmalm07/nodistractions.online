package config

import "testing"

func TestSecretVersionName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "regional secret without version",
			input:  "projects/933874694474/locations/us-east4/secrets/nodistractions-secrets",
			output: "projects/933874694474/locations/us-east4/secrets/nodistractions-secrets/versions/latest",
		},
		{
			name:   "explicit version",
			input:  "projects/example/secrets/app-config/versions/4",
			output: "projects/example/secrets/app-config/versions/4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := secretVersionName(test.input)
			if err != nil {
				t.Fatalf("secretVersionName() error = %v", err)
			}
			if actual != test.output {
				t.Errorf("secretVersionName() = %q, want %q", actual, test.output)
			}
		})
	}
}

func TestSecretLocation(t *testing.T) {
	if actual := secretLocation("projects/933874694474/locations/us-east4/secrets/nodistractions-secrets/versions/latest"); actual != "us-east4" {
		t.Errorf("secretLocation() = %q, want us-east4", actual)
	}
	if actual := secretLocation("projects/example/secrets/app-config/versions/latest"); actual != "" {
		t.Errorf("secretLocation() = %q, want empty string", actual)
	}
}

package docker

import "testing"

func TestImage_String(t *testing.T) {
	tests := []struct {
		name  string
		c     Image
		proxy string
		want  string
	}{
		{
			name:  "Alpine without proxy",
			c:     Alpine_3_18,
			proxy: "",
			want:  "alpine:3.18",
		},
		{
			name:  "Alpine with proxy",
			c:     Alpine_3_18,
			proxy: "${DEPENDENCY_PROXY}/",
			want:  "${DEPENDENCY_PROXY}/alpine:3.18",
		},
		{
			name:  "Unknown (default) without proxy",
			c:     9999,
			proxy: "",
			want:  "alpine:3.18",
		},
		{
			name:  "Unknown (default) with proxy",
			c:     9999,
			proxy: "${DEPENDENCY_PROXY}/",
			want:  "${DEPENDENCY_PROXY}/alpine:3.18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DEPENDENCY_PROXY = tt.proxy
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

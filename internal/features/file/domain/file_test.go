package file_domain

import "testing"

func TestFileStatusIsTerminal(t *testing.T) {
	tests := []struct {
		name   string
		status FileStatus
		want   bool
	}{
		{name: "uploaded is not terminal", status: FileStatusUploaded, want: false},
		{name: "processed is terminal", status: FileStatusProcessed, want: true},
		{name: "failed is terminal", status: FileStatusFailed, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsTerminal(); got != tt.want {
				t.Fatalf("IsTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

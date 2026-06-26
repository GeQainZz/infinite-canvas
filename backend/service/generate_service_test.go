package service

import "testing"

func TestGenerationTypeFromPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "/v1/images/generations", want: "image"},
		{path: "/v1/images/edits", want: "image"},
		{path: "/v1/video/generations", want: "video"},
		{path: "/v1/videos/generations", want: "video"},
		{path: "/v1/videos", want: "video"},
		{path: "/contents/generations/tasks", want: "video"},
		{path: "/v1/audio/speech", want: "audio"},
		{path: "/v1/chat/completions", want: "text"},
		{path: "/v1/responses", want: "text"},
	}

	for _, tt := range tests {
		if got := generationTypeFromPath(tt.path); got != tt.want {
			t.Fatalf("generationTypeFromPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

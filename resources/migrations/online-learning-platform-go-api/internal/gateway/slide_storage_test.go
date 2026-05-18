package gateway

import (
	"testing"

	courseEntity "online-learning-platform-go-api/internal/course/entity"
)

func TestPayloadHasFileReference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		payload courseEntity.PayloadJSON
		want    bool
	}{
		{
			name:    "nil payload",
			payload: nil,
			want:    false,
		},
		{
			name: "empty string ref",
			payload: courseEntity.PayloadJSON{
				"object_key": "   ",
			},
			want: false,
		},
		{
			name: "snake case object key",
			payload: courseEntity.PayloadJSON{
				"object_key": "slides/1/file.docx",
			},
			want: true,
		},
		{
			name: "camel case file src",
			payload: courseEntity.PayloadJSON{
				"fileSrc": "/api/files/slides/1/file.pdf",
			},
			want: true,
		},
		{
			name: "public url",
			payload: courseEntity.PayloadJSON{
				"url": "https://cdn.example.com/slides/1/file.png",
			},
			want: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := payloadHasFileReference(tc.payload)
			if got != tc.want {
				t.Fatalf("payloadHasFileReference() = %v, want %v", got, tc.want)
			}
		})
	}
}

package main

import (
	"testing"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    pb.MediaFormat
		wantErr bool
	}{
		{
			name:    "vinyl lowercase",
			input:   "vinyl",
			want:    pb.MediaFormat_MEDIA_FORMAT_VINYL,
			wantErr: false,
		},
		{
			name:    "vinyl uppercase with spaces",
			input:   "  VINYL  ",
			want:    pb.MediaFormat_MEDIA_FORMAT_VINYL,
			wantErr: false,
		},
		{
			name:    "cd lowercase",
			input:   "cd",
			want:    pb.MediaFormat_MEDIA_FORMAT_CD,
			wantErr: false,
		},
		{
			name:    "cd mixed case",
			input:   "Cd",
			want:    pb.MediaFormat_MEDIA_FORMAT_CD,
			wantErr: false,
		},
		{
			name:    "mixed format",
			input:   "mixed",
			want:    pb.MediaFormat_MEDIA_FORMAT_MIXED,
			wantErr: false,
		},
		{
			name:    "unspecified explicit",
			input:   "unspecified",
			want:    pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			wantErr: false,
		},
		{
			name:    "empty string defaults to unspecified",
			input:   "",
			want:    pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			wantErr: false,
		},
		{
			name:    "whitespace only defaults to unspecified",
			input:   "   ",
			want:    pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			wantErr: false,
		},
		{
			name:    "invalid format returns error",
			input:   "cassette",
			want:    pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

package v1

import (
	"testing"
)

func TestMediaFormatEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		format   MediaFormat
		expected int32
		str      string
	}{
		{
			name:     "Unspecified",
			format:   MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			expected: 0,
			str:      "MEDIA_FORMAT_UNSPECIFIED",
		},
		{
			name:     "Vinyl",
			format:   MediaFormat_MEDIA_FORMAT_VINYL,
			expected: 1,
			str:      "MEDIA_FORMAT_VINYL",
		},
		{
			name:     "CD",
			format:   MediaFormat_MEDIA_FORMAT_CD,
			expected: 2,
			str:      "MEDIA_FORMAT_CD",
		},
		{
			name:     "Mixed",
			format:   MediaFormat_MEDIA_FORMAT_MIXED,
			expected: 3,
			str:      "MEDIA_FORMAT_MIXED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if int32(tc.format) != tc.expected {
				t.Errorf("expected enum value %d, got %d", tc.expected, int32(tc.format))
			}
			if tc.format.String() != tc.str {
				t.Errorf("expected enum string %q, got %q", tc.str, tc.format.String())
			}
		})
	}
}

func TestGenerateDescriptionRequest_GetMediaFormat(t *testing.T) {
	var nilReq *GenerateDescriptionRequest
	if got := nilReq.GetMediaFormat(); got != MediaFormat_MEDIA_FORMAT_UNSPECIFIED {
		t.Errorf("expected nil request GetMediaFormat() to return %v, got %v", MediaFormat_MEDIA_FORMAT_UNSPECIFIED, got)
	}

	defaultReq := &GenerateDescriptionRequest{
		RecordTitle: "Test Album",
		Artist:      "Test Artist",
	}
	if got := defaultReq.GetMediaFormat(); got != MediaFormat_MEDIA_FORMAT_UNSPECIFIED {
		t.Errorf("expected omitted MediaFormat to default to %v, got %v", MediaFormat_MEDIA_FORMAT_UNSPECIFIED, got)
	}

	for _, fmtVal := range []MediaFormat{
		MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
		MediaFormat_MEDIA_FORMAT_VINYL,
		MediaFormat_MEDIA_FORMAT_CD,
		MediaFormat_MEDIA_FORMAT_MIXED,
	} {
		req := &GenerateDescriptionRequest{
			RecordTitle: "Test Album",
			MediaFormat: fmtVal,
		}
		if got := req.GetMediaFormat(); got != fmtVal {
			t.Errorf("expected GetMediaFormat() = %v, got %v", fmtVal, got)
		}
	}
}

package genai_test

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/genai"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	// Build inputs using []byte to guarantee exact byte values survive editor encoding.
	withNull := string([]byte{'h', 'e', 'l', 'l', 'o', 0x00, 'w', 'o', 'r', 'l', 'd'})
	with0x80 := string([]byte{'h', 'e', 'l', 'l', 'o', 0x80, 'w', 'o', 'r', 'l', 'd'})
	with0xFF := string([]byte{'h', 'e', 'l', 'l', 'o', 0xFF, 'w', 'o', 'r', 'l', 'd'})
	withMulti := string([]byte{'a', 0x80, 'b', 0x90, 'c', 0xFE, 'd'})
	// "café Übung" in valid UTF-8
	validUTF8 := "caf\xc3\xa9 \xc3\x9cbung"
	withAds := "(adsbygoogle = window.adsbygoogle || []).push({})"
	withAdsEmbedded := "some text" + withAds + "more text"
	withBoth := string([]byte{'t', 'i', 't', 'l', 'e', 0x80}) + withAds

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean string is unchanged",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "null byte is removed",
			input: withNull,
			want:  "helloworld",
		},
		{
			name:  "0x80 invalid UTF-8 byte is removed",
			input: with0x80,
			want:  "helloworld",
		},
		{
			name:  "0xFF invalid UTF-8 byte is removed",
			input: with0xFF,
			want:  "helloworld",
		},
		{
			name:  "multiple invalid bytes are all removed",
			input: withMulti,
			want:  "abcd",
		},
		{
			name:  "valid UTF-8 multibyte characters are preserved",
			input: validUTF8,
			want:  validUTF8,
		},
		{
			name:  "adsbygoogle artifact is removed",
			input: withAds,
			want:  "",
		},
		{
			name:  "adsbygoogle artifact embedded in text is removed",
			input: withAdsEmbedded,
			want:  "some textmore text",
		},
		{
			name:  "combined invalid bytes and adsbygoogle are both removed",
			input: withBoth,
			want:  "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genai.SanitizeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

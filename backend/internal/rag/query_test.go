package rag

import "testing"

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		userProfile string
		want        string
	}{
		{
			name:  "query without profile",
			query: "four hundred meter endurance",
			want:  "four hundred meter endurance",
		},
		{
			name:        "query with profile",
			query:       "four hundred meter endurance",
			userProfile: "Benutzerprofil Präferenzen:\n- Erfahrungslevel: Fortgeschritten",
			want:        "four hundred meter endurance\n\nBenutzerprofil Präferenzen:\n- Erfahrungslevel: Fortgeschritten",
		},
		{
			name:        "profile without query",
			userProfile: "Benutzerprofil Präferenzen:\n- Erfahrungslevel: Anfänger",
			want:        "Benutzerprofil Präferenzen:\n- Erfahrungslevel: Anfänger",
		},
		{
			name:        "whitespace is ignored",
			query:       "  endurance  ",
			userProfile: "\n profile \n",
			want:        "endurance\n\nprofile",
		},
		{
			name:        "empty values",
			query:       "  ",
			userProfile: "\n",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSearchQuery(tt.query, tt.userProfile); got != tt.want {
				t.Errorf("buildSearchQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}

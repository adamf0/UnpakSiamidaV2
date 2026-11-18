package domain

import "testing"

func TestProkerV3Default_DisplayFakultasUnit(t *testing.T) {
	cases := []struct {
		name string
		p    ProkerV3Default
		want string
	}{
		{
			name: "prodi with jenjang appends (S1)",
			p:    ProkerV3Default{NamaFakultasUnit: "ILMU KOMPUTER", TypeFakultasUnit: "prodi", JenjangFakultasUnit: "s1"},
			want: "ILMU KOMPUTER (S1)",
		},
		{
			name: "prodi already containing (S1) does not duplicate",
			p:    ProkerV3Default{NamaFakultasUnit: "ILMU KOMPUTER (S1)", TypeFakultasUnit: "prodi", JenjangFakultasUnit: "S1"},
			want: "ILMU KOMPUTER (S1)",
		},
		{
			name: "non-prodi stays plain name",
			p:    ProkerV3Default{NamaFakultasUnit: "MIPA", TypeFakultasUnit: "fakultas", JenjangFakultasUnit: ""},
			want: "MIPA",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.DisplayFakultasUnit(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

package infrastructure

import "testing"

func TestComputeTitleSimilarity(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "exact match with case and punctuation",
			a:    "Lisensi dan Pengembangan Infrastruktur IT Aset Inventaris (YPS)",
			b:    "lisensi dan pengembangan infrastruktur it aset inventaris yps",
			want: true,
		},
		{
			name: "minor typo",
			a:    "Lisensi dan Pengembangan Infrastrukter IT Aset Inventaris",
			b:    "Lisensi dan Pengembangan Infrastruktur IT Aset Inventaris",
			want: true,
		},
		{
			name: "numbered series should not match",
			a:    "Rapat Evaluasi Kinerja 1",
			b:    "Rapat Evaluasi Kinerja 2",
			want: false,
		},
		{
			name: "different activities",
			a:    "Pengadaan Komputer Laboratorium",
			b:    "Workshop Penyusunan Kurikulum OBE",
			want: false,
		},
		{
			name: "empty title",
			a:    "",
			b:    "Pelatihan Audit",
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsSimilarProgramTitle(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("IsSimilarProgramTitle(%q, %q) = %v (score %.3f), want %v", tc.a, tc.b, got, ComputeTitleSimilarity(tc.a, tc.b), tc.want)
			}
		})
	}
}

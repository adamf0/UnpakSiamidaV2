package application

type UpdateTemplateDokumenTambahanCommand struct {
    Uuid     	 	string
	Tahun     	    string
    JenisFile       string //uuid
    Pertanyaan      string
    Klasifikasi     string //minor, major
    Kategori     	string
    Tugas     	    string //auditor1, auditor2
}

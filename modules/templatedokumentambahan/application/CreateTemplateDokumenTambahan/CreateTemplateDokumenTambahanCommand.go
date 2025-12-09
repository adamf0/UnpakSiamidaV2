package application

type CreateTemplateDokumenTambahanCommand struct {
    Tahun     	    string
    JenisFile       string //uuid
    Pertanyaan      string
    Klasifikasi     string //minor, major
    Kategori     	string
    Tugas     	    string //auditor1, auditor2
}

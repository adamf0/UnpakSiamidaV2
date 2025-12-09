package application

type CreateTemplateRenstraCommand struct {
    Tahun     	    string
    Indikator       string //uuid
    IsPertanyaan    string
    FakultasUnit    string //uuid
    Kategori     	string
    Klasifikasi     string //minor, major
    Satuan     	    *string //rule aktif jika bukan nil
    Target     	    *string //jika target && (targetmin || target max) kosong maka required; jika target tidak kosong maka skip targetmin targetmax
    TargetMin     	*string //kebalikan dari rule target
    TargetMax     	*string //kebalikan rule target
    Tugas     	    string //auditor1, auditor2
}

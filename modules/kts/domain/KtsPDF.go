package domain

import "html/template"

type Auditee struct {
	Nama string `json:"nama"`
}

type AccAuditor struct {
	Nama string `json:"nama"`
	Tgl  string `json:"tgl,omitempty"`
}

type CloseData struct {
	WmmUpmfUpmps   string `json:"wmm_upmf_upmps"`
	TanggalClosing string `json:"tanggal_closing"`
}

type DetailPoint struct {
	NamaStandarPointing string         `json:"nama_standar_pointing"`
	Kesimpulan          string         `json:"kesimpulan"`
	Counting            map[string]int `json:"counting,omitempty"` // minor/mayor
	Point               []DetailPoint  `json:"point,omitempty"`    // nested point
}

type KtsPDF struct {
	Nomor                     string        `json:"nomor"`
	Tanggal                   string        `json:"tanggal"`
	Auditee                   Auditee       `json:"auditee"`
	Auditor1                  string        `json:"auditor1"`
	AccAuditor1               AccAuditor    `json:"acc_auditor1"`
	AccAuditor1Final          AccAuditor    `json:"acc_auditor1_final"`
	P                         template.HTML `json:"p"`
	L                         template.HTML `json:"l"`
	O                         template.HTML `json:"o"`
	R                         template.HTML `json:"r"`
	AkarMasalah               template.HTML `json:"akar_masalah"`
	TindakanKoreksi           template.HTML `json:"tindakan_koreksi"`
	Referensi                 string        `json:"referensi"`
	HasilTemuan               string        `json:"hasil_temuan"`
	Detail                    []DetailPoint `json:"detail"`
	TindakanPerbaikan         string        `json:"tindakan_perbaikan"`
	TanggalPenyelesaian       string        `json:"tanggal_penyelesaian"`
	TinjauanTindakanPerbaikan string        `json:"tinjauan_tindakan_perbaikan"`
	Close                     CloseData     `json:"close"`
	QRAccAuditor1             string        `json:"qr_acc_auditor1"`
	QRAuditee                 string        `json:"qr_auditee"`
	QRAccAuditor1Final        string        `json:"qr_acc_auditor1_final"`
	QRClose                   string        `json:"qr_close"`
}

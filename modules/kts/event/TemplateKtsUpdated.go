package event

const ktsUpdatedTemplate = `
✏️ <b>Perubahan Data KTS</b>

<b>KTS</b>
• UUID : <code>{{.UUID}}</code>
• Target : <b>{{.Target}}</b>
• Tahun : <b>{{.Tahun}}</b>
• Status : <b>{{.Status}}</b>
• Nomor Laporan : <b>{{.NomorLaporan}}</b>

<b>Konteks Renstra</b>
• Standar   : {{.Standar}}
• Indikator : {{.Indikator}}

<b>Konteks Dokumen</b>
• Pertanyaan : {{.Pertanyaan}}
• Jenis File : {{.JenisFile}}

<b>Ketidaksesuaian</b>
<pre>
P : {{.P}}
L : {{.L}}
O : {{.O}}
R : {{.R}}
</pre>

<b>Tindakan</b>
• Akar Masalah : {{.AkarMasalah}}
• Tindakan Koreksi : {{.TindakanKoreksi}}

<b>Respon Auditee</b>
<pre>
Status                     : {{.StatusAcc}}
Keterangan                 : {{.Keterangan}}
Tindakan Perbaikan         : {{.TindakanPerbaikan}}
Tinjauan Tindakan Perbaikan: {{.Tinjauan}}
Tanggal Closing            : {{.TanggalClosing}}
</pre>

<b>Respon Auditor</b>
<pre>
• Tanggal Penyelesaian : {{.TanggalPenyelesaian}}
• Tanggal Closing Final : {{.TanggalClosingFinal}}
• WmmUpmfUpmps : {{.Wmm}}
</pre>

<b>Waktu</b>
• Terjadi : {{.Terjadi}}
`

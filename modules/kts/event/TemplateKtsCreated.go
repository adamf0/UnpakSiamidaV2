package event

const ktsCreatedTemplate = `
🆕 <b>KTS Baru Dibuat</b>

<b>Identitas</b>
• KTS UUID : <code>{{.KtsUUID}}</code>

<b>Target</b>
• Target : <b>{{.Target}}</b>
• Tahun  : <b>{{.Tahun}}</b>
• Status : <b>{{.Status}}</b>

<b>Konteks Renstra</b>
• Standar   : {{.Standar}}
• Indikator : {{.Indikator}}

<b>Konteks Dokumen</b>
• Pertanyaan : {{.Pertanyaan}}
• Jenis File : {{.JenisFile}}

<b>Template</b>
• Template Renstra : {{.TemplateRenstra}}
• Template Dokumen : {{.TemplateDokumen}}

<b>Metadata</b>
• Terjadi : {{.OccurredOn}}
`

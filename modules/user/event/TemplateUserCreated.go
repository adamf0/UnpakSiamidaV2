package event

const userCreatedTemplate = `
👤 <b>Data User</b>

<b>Identitas</b>
• Username : <b>{{.Username}}</b>
• Nama     : <b>{{.Name}}</b>
• Email    : <b>{{.Email}}</b>

<b>Akses</b>
• Level : <b>{{.Level}}</b>

<b>Tugaskan</b>
• Fakultas / Unit : <b>{{.FakultasUnit}}</b>
• Tipe  : <b>{{.Tipe}}</b>

<b>Keamanan</b>
• Password : <code>{{.Password}}</code>

<b>Waktu</b>
• Terjadi : {{.Terjadi}}
`

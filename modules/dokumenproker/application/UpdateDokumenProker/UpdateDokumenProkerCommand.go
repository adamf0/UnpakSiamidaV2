package application

type UpdateDokumenProkerCommand struct {
	Uuid            string
	FakultasUuid    string
	MataProgramUuid string
	JenisDokumen    string
	File            string
	Status          string
	Catatan         *string
}

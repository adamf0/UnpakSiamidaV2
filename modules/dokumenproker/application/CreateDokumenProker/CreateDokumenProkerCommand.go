package application

type CreateDokumenProkerCommand struct {
	FakultasUuid    string
	MataProgramUuid string
	JenisDokumen    string
	File            string
	Status          string
	Catatan         *string
}

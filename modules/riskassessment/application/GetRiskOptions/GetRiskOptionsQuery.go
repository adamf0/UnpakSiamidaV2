package GetRiskOptions

type GetRiskOptionsQuery struct {
	FakultasUnit string
	Tahun        string
}

type GetRiskCriteriaQuery struct {
	Tipe string
}

type GetRiskLevelsQuery struct{}

type GetRiskCategoriesQuery struct{}

type GetRiskStatusesQuery struct{}

type GetRiskProgramsQuery struct {
	FakultasUnit string
	Tahun        string
	Search       string
}

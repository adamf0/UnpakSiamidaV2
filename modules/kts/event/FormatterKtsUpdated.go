package event

import (
	helper "UnpakSiamida/common/helper"
	"bytes"
	"html/template"
)

type KtsUpdatedView struct {
	UUID                string
	Target              string
	Tahun               string
	Status              string
	NomorLaporan        string
	Standar             string
	Indikator           string
	Pertanyaan          string
	JenisFile           string
	P                   string
	L                   string
	O                   string
	R                   string
	AkarMasalah         string
	TindakanKoreksi     string
	StatusAcc           string
	Keterangan          string
	TindakanPerbaikan   string
	Tinjauan            string
	TanggalClosing      string
	TanggalPenyelesaian string
	TanggalClosingFinal string
	Wmm                 string
	Terjadi             string
}

func RenderKtsUpdatedTemplate(e KtsUpdatedEvent) string {
	ctx := helper.DateContext{}
	ctx.SetStrategy(helper.IndonesianDateFormatter{})

	factory := helper.DateLayoutFirstFactory{}
	TanggalClosing := helper.NewDateChain(*e.TanggalClosing).
		UseParseStrategy(factory.CreateParser()).
		Parse().
		UseFormatStrategy(factory.CreateFormatter()).
		FormatString()

	TanggalPenyelesaian := helper.NewDateChain(*e.TanggalPenyelesaian).
		UseParseStrategy(factory.CreateParser()).
		Parse().
		UseFormatStrategy(factory.CreateFormatter()).
		FormatString()

	TanggalClosingFinal := helper.NewDateChain(*e.TanggalClosingFinal).
		UseParseStrategy(factory.CreateParser()).
		Parse().
		UseFormatStrategy(factory.CreateFormatter()).
		FormatString()

	view := KtsUpdatedView{
		UUID:                e.KtsUUID.String(),
		Target:              helper.StringHtmlValue(e.Target),
		Tahun:               helper.StringHtmlValue(e.Tahun),
		Status:              helper.StringHtmlValue(e.Status),
		NomorLaporan:        helper.StringHtmlValue(e.NomorLaporan),
		Standar:             helper.StringHtmlValue(e.Standar),
		Indikator:           helper.StringHtmlValue(e.Indikator),
		Pertanyaan:          helper.StringHtmlValue(e.Pertanyaan),
		JenisFile:           helper.StringHtmlValue(e.JenisFile),
		P:                   helper.StringHtmlValue(e.KetidaksesuaianP),
		L:                   helper.StringHtmlValue(e.KetidaksesuaianL),
		O:                   helper.StringHtmlValue(e.KetidaksesuaianO),
		R:                   helper.StringHtmlValue(e.KetidaksesuaianR),
		AkarMasalah:         helper.StringHtmlValue(e.AkarMasalah),
		TindakanKoreksi:     helper.StringHtmlValue(e.TindakanKoreksi),
		StatusAcc:           helper.Status(e.StatusAccAuditee),
		Keterangan:          helper.StringHtmlValue(e.KeteranganTolak),
		TindakanPerbaikan:   helper.StringHtmlValue(e.TindakanPerbaikan),
		Tinjauan:            helper.StringHtmlValue(e.TinjauanTindakanPerbaikan),
		TanggalClosing:      TanggalClosing,
		TanggalPenyelesaian: TanggalPenyelesaian,
		TanggalClosingFinal: TanggalClosingFinal,
		Wmm:                 helper.StringHtmlValue(e.WmmUpmfUpmps),
		Terjadi:             ctx.FormatWithTime(e.OccurredOn),
	}

	tpl := template.Must(template.New("kts").Parse(ktsUpdatedTemplate))

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, view)

	return buf.String()
}

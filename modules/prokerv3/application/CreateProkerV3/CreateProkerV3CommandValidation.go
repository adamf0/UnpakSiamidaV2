package CreateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"strings"
)

func (c *CreateProkerV3Command) Validate() error {
	if strings.TrimSpace(c.NamaProgram) == "" {
		return domain.NamaProgramRequired()
	}
	return nil
}

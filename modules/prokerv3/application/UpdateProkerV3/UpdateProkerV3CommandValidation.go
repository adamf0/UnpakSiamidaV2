package UpdateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"strings"

	"github.com/google/uuid"
)

func (c *UpdateProkerV3Command) Validate() error {
	if _, err := uuid.Parse(strings.TrimSpace(c.UUID)); err != nil {
		return domain.ProkerNotFound()
	}
	if strings.TrimSpace(c.NamaProgram) == "" {
		return domain.NamaProgramRequired()
	}
	return nil
}

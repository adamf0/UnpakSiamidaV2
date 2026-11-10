package DeleteProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"strings"

	"github.com/google/uuid"
)

func (c *DeleteProkerV3Command) Validate() error {
	if _, err := uuid.Parse(strings.TrimSpace(c.UUID)); err != nil {
		return domain.ProkerNotFound()
	}
	if c.SoftDelete && strings.TrimSpace(c.Catatan) == "" {
		return domain.DeleteReasonRequired()
	}
	return nil
}

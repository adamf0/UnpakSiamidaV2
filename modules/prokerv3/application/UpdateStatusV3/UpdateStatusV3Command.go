package UpdateStatusV3

type UpdateStatusV3Command struct {
	UUID             string `json:"uuid" form:"uuid"`
	Status           string `json:"status" form:"status"`
	Role             string `json:"role" form:"role"`
	CatatanPenolakan string `json:"catatan_penolakan" form:"catatan_penolakan"`
	Alasan           string `json:"alasan" form:"alasan"`
	Catatan          string `json:"catatan" form:"catatan"`
	Level            string `json:"level" form:"level"`
	UserFU           string `json:"user_fu" form:"user_fu"`
	UserID           *int   `json:"user_id" form:"user_id"`
	UserName         string `json:"user_name" form:"user_name"`
}

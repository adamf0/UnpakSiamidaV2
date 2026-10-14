package application

type MergeUserCommand struct {
	PrimaryUUID    string   `json:"primary_uuid"`
	SecondaryUUIDs []string `json:"secondary_uuids"`
	PrimaryID      *uint    `json:"primary_id"`
	SecondaryIDs   []uint   `json:"secondary_ids"`
}

package DeleteProkerV3

type DeleteProkerV3Command struct {
	UUID          string `json:"uuid"`
	SoftDelete    bool   `json:"soft_delete"`
	Catatan       string `json:"catatan"`
	Role          string `json:"role"`
	ActorID       *int   `json:"actor_id"`
	ActorName     string `json:"actor_name"`
	ActorUsername string `json:"actor_username"`
}

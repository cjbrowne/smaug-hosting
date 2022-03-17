package containers

type ContainerStatus struct {
	Up    bool   `json:"up"`
	State string `json:"state"`
}

type Container struct {
	Id        int64           `json:"id"`
	Name      string          `json:"name"`
	Tier      int             `json:"tier"`
	Software  string          `json:"software"`
	UserId    int64           `json:"-" db:"user_id"`
	LastError string          `json:"last_error"`
	Status    ContainerStatus `json:"status"`
	IP        string          `json:"ip"`
	Port      uint32             `json:"port"`
}

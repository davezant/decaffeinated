package api

func (c *Client) GetWatcher() (*WatcherSnapshot, error) {
	var snapshot WatcherSnapshot
	err := c.get("/watcher/", &snapshot)
	return &snapshot, err
}

type WatcherSnapshot struct {
	ProcessesSnapshot struct {
		Processes []string `json:"Processes"`
	} `json:"ProcessesSnapshot"`

	ActivitiesUp struct {
		Active   []Activity `json:"Active"`
		Inactive []Activity `json:"Inactive"`
	} `json:"ActivitiesUp"`

	ActiveSession struct {
		UserID    string `json:"UserID"`
		LoginTime string `json:"LoginTime"`
		Limit     int64  `json:"Limit"`
		IsMinor   bool   `json:"IsMinor"`
	} `json:"ActiveSession"`

	ServiceStartTime     string `json:"ServiceStartTime"`
	SessionExecutionTime string `json:"SessionExecutionTime"`
}

type Activity struct {
	Name                 string `json:"Name"`
	ExecutionBinary      string `json:"ExecutionBinary"`
	IsUp                 bool   `json:"IsUp"`
	Limit                int64  `json:"Limit"`
	DisplayExecutionTime string `json:"DisplayExecutionTime"`
	IsCounting           bool   `json:"IsCounting"`
}

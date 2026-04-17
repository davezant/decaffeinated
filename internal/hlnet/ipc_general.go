package hlnet

const (
	DefaultLinuxSockPath = `./sock`
	DefaultWindowsPipePath = `\\.\pipe\decaffeinated`
)

type IPCConfig struct {
	Path string
}

type IPCCommand struct {
	Action           string `json:"action"`
	AppNames          []string `json:"app_names"`
	Category		 string `json:"category"`
	TimeLimitSeconds int64  `json:"time_limit_seconds,omitempty"`
	IsBlocked        bool   `json:"is_blocked,omitempty"`
	CustomTimestamps []CustomTimestamp `json:"rule,omitempty"`
}

type CustomTimestamp struct {
	Timestamp float32 `json:"timestamp"`
	Callback  string  `json:"callback"`
}

// Heartbeat for external User control.
type Heartbeat struct {
	UserUUID string `json:"UUID"`
	TimeLimitSeconds int64 `json:"time_limit_seconds,omitempty"`
}

// IPCResponse is returned by the WatchDog IPC endpoint.
type IPCResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CommandBundle struct {
	Commands []IPCCommand
}

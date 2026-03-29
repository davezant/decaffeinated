package dwatchdog

type NetConfig struct {
	BlockedIPS map[string]bool
	BlockedHostnames map[string]bool
}


package dwatchdog

type NetConfig struct {
	BlockedIPS map[string]bool
	BlockedHostnames map[string]bool
	Host string
	Port string
}

type NetManager interface {
	NewProxy(host, port string) error
	StartProxy() error
	StopProxy() error
	BlockIP(ip string) error
	BlockHostname(hostname string) error
}

type NetProxy struct {

}
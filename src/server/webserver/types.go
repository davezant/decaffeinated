package webserver

const DefaultURL = "localhost"
const DefaultPort = "1337" // Latte

type WebserverConfigs struct {
	URL  string
	PORT int
}

type ProxyConfigs struct {
	URL  string
	PORT int
}

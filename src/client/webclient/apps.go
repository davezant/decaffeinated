package webclient

import "time"

func (c *Client) GetApps() ([]App, error) {
	var apps []App
	err := c.get("/apps/", &apps)
	return apps, err
}

func (c *Client) GetApp(name string) (*App, error) {
	var app App
	err := c.get("/apps/"+name, &app)
	return &app, err
}

func (c *Client) CreateApp(name, binary, path, prefix, suffix string, limit time.Duration, minors bool) (*App, error) {
	payload := map[string]interface{}{
		"name":            name,
		"binary":          binary,
		"path":            path,
		"limit":           limit,
		"command_prefix":  prefix,
		"command_suffix":  suffix,
		"can_minors_play": minors,
	}
	var app App
	err := c.post("/apps/", payload, &app)
	return &app, err
}

func (c *Client) DeleteApp(name string) error {
	var resp map[string]string
	return c.delete("/apps/"+name, &resp)
}

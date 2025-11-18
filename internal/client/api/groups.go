package api

func (c *Client) GetGroups() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := c.get("/groups", &data)
	return data, err
}

func (c *Client) GetGroup(name string) (*Group, error) {
	var group Group
	err := c.get("/groups/"+name, &group)
	return &group, err
}

func (c *Client) CreateGroup(name string) (*Group, error) {
	var resp struct {
		Status string `json:"status"`
		Group  Group  `json:"group"`
	}
	err := c.put("/groups", map[string]string{"name": name}, &resp)
	return &resp.Group, err
}

func (c *Client) DeleteGroup(name string) error {
	var resp map[string]string
	return c.delete("/groups/"+name, &resp)
}

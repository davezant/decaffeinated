package webclient

func (c *Client) GetSession() (*Session, error) {
	var s Session
	err := c.get("/session/", &s)
	return &s, err
}

func (c *Client) Login(user, pass string) (*Session, error) {
	var s Session
	err := c.post("/login/", map[string]string{
		"name":     user,
		"password": pass,
	}, &s)
	return &s, err
}

func (c *Client) Logout() error {
	var resp map[string]string
	return c.get("/logout/", &resp)
}

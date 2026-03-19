package api

// SendChatMessage sends a chat message to a running meeting.
func (c *Client) SendChatMessage(req *SendChatMessageRequest) (*SendChatMessageResponse, error) {
	body, err := c.doGet("sendChatMessage", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &SendChatMessageResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetJoinUrl gets a join URL from a session token.
func (c *Client) GetJoinUrl(req *GetJoinUrlRequest) (*GetJoinUrlResponse, error) {
	body, err := c.doGet("getJoinUrl", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &GetJoinUrlResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

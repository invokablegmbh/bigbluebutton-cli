package api

// IsMeetingRunning checks if a meeting is currently running.
func (c *Client) IsMeetingRunning(req *IsMeetingRunningRequest) (*IsMeetingRunningResponse, error) {
	body, err := c.doGet("isMeetingRunning", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &IsMeetingRunningResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMeetingInfo returns detailed information about a meeting.
func (c *Client) GetMeetingInfo(req *GetMeetingInfoRequest) (*GetMeetingInfoResponse, error) {
	body, err := c.doGet("getMeetingInfo", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &GetMeetingInfoResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMeetings returns a list of all running meetings.
func (c *Client) GetMeetings() (*GetMeetingsResponse, error) {
	body, err := c.doGet("getMeetings", (&GetMeetingsRequest{}).ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &GetMeetingsResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

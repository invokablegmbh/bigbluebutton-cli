package api

import (
	"strings"
)

// Create creates a new BigBlueButton meeting.
func (c *Client) Create(req *CreateRequest) (*CreateResponse, error) {
	body, err := c.doGet("create", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &CreateResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// Join generates a join URL for a BigBlueButton meeting.
// When redirect is set to "false", it returns meeting details instead of redirecting.
func (c *Client) Join(req *JoinRequest) (*JoinResponse, error) {
	// For join, we typically want the URL, not to follow the redirect
	if req.Redirect == nil {
		req.Redirect = StringPtr("false")
	}
	body, err := c.doGet("join", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &JoinResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// JoinURL returns the direct join URL (for use in browsers) without making an API call.
func (c *Client) JoinURL(req *JoinRequest) string {
	params := req.ToValues()
	qs := params.Encode()
	checksum := Checksum("join", qs, c.Secret)
	return c.BaseURL + "/join?" + qs + "&checksum=" + checksum
}

// End terminates a running BigBlueButton meeting.
func (c *Client) End(req *EndRequest) (*EndResponse, error) {
	body, err := c.doGet("end", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &EndResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// InsertDocument inserts a presentation document into a running meeting.
func (c *Client) InsertDocument(req *InsertDocumentRequest) (*InsertDocumentResponse, error) {
	body, err := c.doPost("insertDocument", req.ToValues(), "application/xml", strings.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &InsertDocumentResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

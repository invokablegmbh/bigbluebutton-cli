package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// GetRecordings returns a list of recordings.
func (c *Client) GetRecordings(req *GetRecordingsRequest) (*GetRecordingsResponse, error) {
	body, err := c.doGet("getRecordings", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &GetRecordingsResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// PublishRecordings publishes or unpublishes recordings.
func (c *Client) PublishRecordings(req *PublishRecordingsRequest) (*PublishRecordingsResponse, error) {
	body, err := c.doGet("publishRecordings", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &PublishRecordingsResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteRecordings deletes recordings.
func (c *Client) DeleteRecordings(req *DeleteRecordingsRequest) (*DeleteRecordingsResponse, error) {
	body, err := c.doGet("deleteRecordings", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &DeleteRecordingsResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateRecordings updates recording metadata.
func (c *Client) UpdateRecordings(req *UpdateRecordingsRequest) (*UpdateRecordingsResponse, error) {
	body, err := c.doGet("updateRecordings", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseXMLResponse(body, &UpdateRecordingsResponse{})
	if err != nil {
		return nil, err
	}
	if err := resp.ToError(); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRecordingTextTracks returns text tracks for a recording (JSON response).
func (c *Client) GetRecordingTextTracks(req *GetRecordingTextTracksRequest) (*GetRecordingTextTracksResponse, error) {
	body, err := c.doGet("getRecordingTextTracks", req.ToValues())
	if err != nil {
		return nil, err
	}
	resp, err := parseJSONResponse(body, &GetRecordingTextTracksResponse{})
	if err != nil {
		return nil, err
	}
	if resp.Response.ReturnCode != "SUCCESS" {
		return nil, &APIError{
			ReturnCode: resp.Response.ReturnCode,
			MessageKey: resp.Response.MessageKey,
			Message:    resp.Response.Message,
		}
	}
	return resp, nil
}

// PutRecordingTextTrack uploads a subtitle file for a recording.
func (c *Client) PutRecordingTextTrack(req *PutRecordingTextTrackRequest) (*PutRecordingTextTrackResponse, error) {
	// Build multipart form body
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("opening subtitle file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(req.FilePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copying file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	body, err := c.doPost("putRecordingTextTrack", req.ToValues(), writer.FormDataContentType(), &buf)
	if err != nil {
		return nil, err
	}
	resp, err := parseJSONResponse(body, &PutRecordingTextTrackResponse{})
	if err != nil {
		return nil, err
	}
	if resp.Response.ReturnCode != "SUCCESS" {
		return nil, &APIError{
			ReturnCode: resp.Response.ReturnCode,
			MessageKey: resp.Response.MessageKey,
			Message:    resp.Response.Message,
		}
	}
	return resp, nil
}

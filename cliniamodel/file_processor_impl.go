package cliniamodel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/clinia/models-client-go/cliniamodel/requesterhttp/filesvcclient"
	"github.com/clinia/x/errorx"
)

type fileProcessor struct {
	client *filesvcclient.Client
}

var _ FileProcessor = (*fileProcessor)(nil)

func NewFileProcessor(baseURL string, opts ...filesvcclient.ClientOption) (FileProcessor, error) {
	c, err := filesvcclient.NewClient(baseURL, opts...)
	if err != nil {
		return nil, err
	}

	return fileProcessor{client: c}, nil
}

func (c fileProcessor) SplitPDFToImages(
	ctx context.Context,
	body PDFSplitRequest,
	reqEditors ...filesvcclient.RequestEditorFn,
) (zipBytes []byte, err error) {
	{
		if len(body.PDF) == 0 {
			return nil, errorx.InvalidArgumentError("pdf bytes are empty")
		}
		if body.Filename == "" {
			// Arbitrary
			body.Filename = "file.pdf"
		}

		// Build multipart form
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)

		fw, err := w.CreateFormFile("file", body.Filename)
		if err != nil {
			return nil, errorx.InternalErrorf("create form file: %v", err)
		}
		if _, err := fw.Write(body.PDF); err != nil {
			return nil, errorx.InternalErrorf("write pdf to form file: %v", err)
		}

		if body.DPI != nil {
			err = w.WriteField("dpi", strconv.Itoa(*body.DPI))
			if err != nil {
				return nil, errorx.InternalErrorf("write dpi to form field: %v", err)
			}
		}

		if err := w.Close(); err != nil {
			return nil, errorx.InternalErrorf("close multipart writer: %v", err)
		}

		// Use the generated request helper (no manual URL building)
		resp, err := c.client.SplitToImagesWithBody(ctx, w.FormDataContentType(), &buf, reqEditors...)
		if err != nil {
			return nil, errorx.InternalErrorf("request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errorx.InternalErrorf("read response: %v", err)
		}

		// Success (prefer content-type check, but accept 200 regardless)
		if resp.StatusCode == http.StatusOK {
			if resp.Header.Get("Content-Type") == "application/zip" {
				return body, nil
			}
			return nil, errorx.InternalErrorf("expected zip content-type, got %q", resp.Header.Get("Content-Type"))
		}

		// Try to parse structured Fastapi JSON error: {"detail": {"code": ..., "type": "...", "message": "..."}}
		var ve struct {
			Detail struct {
				Code    int    `json:"code"`
				Type    string `json:"type"`
				Message string `json:"message"`
			} `json:"detail"`
		}
		if json.Unmarshal(body, &ve) == nil && ve.Detail.Type != "" && ve.Detail.Message != "" {
			// Normalize to "[TYPE]" message" and let errorx parse it
			if ce, perr := errorx.NewCliniaErrorFromMessage(fmt.Sprintf("[%s] %s", ve.Detail.Type, ve.Detail.Message)); perr == nil {
				return nil, *ce
			}
		}

		// Fallback: maybe the body already is in "[TYPE] message" format
		if ce, perr := errorx.NewCliniaErrorFromMessage(string(body)); perr == nil {
			return nil, *ce
		}

		return nil, errorx.InternalErrorf("unexpected status %d, body: %s", resp.StatusCode, string(body))
	}
}

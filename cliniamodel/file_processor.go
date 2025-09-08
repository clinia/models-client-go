package cliniamodel

import (
	"context"

	"github.com/clinia/models-client-go/cliniamodel/requesterhttp/filesvcclient"
)

// SplitOpts are optional form fields for /v1/pdf/split/images.
type PDFSplitRequest struct {
	Filename string // optional; defaults to "file.pdf"
	DPI      *int   // optional (72..600 typical); server default applies if nil
	PDF      []byte
}

// FileProcessorClient is the high-level interface your code should depend on.
// Implement it by wrapping the generated filesvc client (or a hand-written http client).
type FileProcessor interface {
	// SplitToImagesBytes posts a PDF (as raw bytes) and returns the ZIP payload bytes on success.
	//
	// Returns:
	//   - zipBytes on success (HTTP 200)
	//   - error conveying provider/transport issues
	SplitPDFToImages(
		ctx context.Context,
		body PDFSplitRequest,
		reqEditors ...filesvcclient.RequestEditorFn,
	) (zipBytes []byte, err error)
}

package response

import (
	"fmt"
	"github.com/kiquetal/learn-http-protocol/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Response struct {
	StatusCode StatusCode
	Headers    headers.Header
	Body       []byte
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusText := getStatusText(statusCode)
	if statusText == "" {
		return fmt.Errorf("unknown status code: %d", statusCode)
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, statusText)
	return err
}

func getStatusText(statusCode StatusCode) string {
	switch statusCode {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}

func GetDefaultHeaders(contentLen int) headers.Header {

	return headers.Header{
		"Content-Length": fmt.Sprintf("%d", contentLen),
		"Content-Type":   "text/plain",
		"Connection":     "close",
	}
}

func WriteHeaders(w io.Writer, headers headers.Header) error {
	for key, value := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", key, value); err != nil {
			return fmt.Errorf("error writing header %s: %w", key, err)
		}
	}
	_, err := w.Write([]byte("\r\n")) // End of headers
	return err
}

type WriterStatus int

const (
	WriterStatusInitialized WriterStatus = iota
	WriterStatusWritingHeaders
	WriterStatusWritingBody
	WriterStatusDone
	WriterStatusError
)

type Writer struct {
	io.Writer   // Underlying writer to write the response
	writeStatus WriterStatus
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	w.writeStatus = WriterStatusInitialized
	statusText := getStatusText(statusCode)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	if _, err := fmt.Fprint(w, statusLine); err != nil {
		w.writeStatus = WriterStatusError
		return fmt.Errorf("error writing status line: %w", err)
	}
	w.writeStatus = WriterStatusWritingHeaders
}

func (w *Writer) WriteHeaders(headers headers.Header) error {
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if len(body) == 0 {
		return nil // No body to write
	}
	_, err := fmt.Fprintf(w, "%s", body)
	return err
}

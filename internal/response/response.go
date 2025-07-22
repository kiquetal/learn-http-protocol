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
		"Content-Type":   "text/html",
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
	WriteStatus WriterStatus
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	w.WriteStatus = WriterStatusInitialized
	statusText := getStatusText(statusCode)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	if _, err := fmt.Fprint(w, statusLine); err != nil {
		w.WriteStatus = WriterStatusError
		return fmt.Errorf("error writing status line: %w", err)
	}
	w.WriteStatus = WriterStatusWritingHeaders
	return nil

}

func (w *Writer) WriteHeaders(headers headers.Header) error {
	if w.WriteStatus != WriterStatusWritingHeaders {
		return fmt.Errorf("cannot write headers in current state: %v", w.WriteStatus)
	}

	if err := WriteHeaders(w, headers); err != nil {
		w.WriteStatus = WriterStatusError
		return fmt.Errorf("error writing headers: %w", err)
	}
	w.WriteStatus = WriterStatusWritingBody
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {

	if w.WriteStatus != WriterStatusWritingBody {
		return 0, fmt.Errorf("cannot write body in current state: %v", w.WriteStatus)
	}
	//need to add the header with the content length
	contentLength := len(body)
	header := headers.Header{
		"Content-Length": fmt.Sprintf("%d", contentLength),
	}
	if err := WriteHeaders(w, header); err != nil {
		n, err := w.Writer.Write(body)
		if err != nil {
			w.WriteStatus = WriterStatusError
			return n, fmt.Errorf("error writing body: %w", err)
		}

	}
	n, err := w.Write(body)
	if err != nil {
		return 0, err
	}
	w.WriteStatus = WriterStatusDone
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	// write chunked response, modify the header
	if w.WriteStatus != WriterStatusWritingBody {
		return 0, fmt.Errorf("cannot write chunked body in current state: %v", w.WriteStatus)
	}
	if len(p) == 0 {
		return 0, nil // No data to write
	}
	chunkSize := fmt.Sprintf("%x\r\n", len(p))
	if _, err := w.Write([]byte(chunkSize)); err != nil {
		w.WriteStatus = WriterStatusError
		return 0, fmt.Errorf("error writing chunk size: %w", err)
	}
	n, err := w.Write(p)
	if err != nil {
		w.WriteStatus = WriterStatusError
		return n, fmt.Errorf("error writing chunked body: %w", err)
	}
	if _, err := w.Write([]byte("\r\n")); err != nil {
		w.WriteStatus = WriterStatusError
		return n, fmt.Errorf("error writing chunked body end: %w", err)
	}
	return n, nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {

	if w.WriteStatus != WriterStatusWritingBody {
		return 0, fmt.Errorf("cannot write chunked body done in current state: %v", w.WriteStatus)
		//write the last chunk with size 0
	}

	if _, err := w.Write([]byte("0\r\n\r\n")); err != nil {
		w.WriteStatus = WriterStatusError
		return 0, fmt.Errorf("error writing chunked body done: %w", err)
	}
	w.WriteStatus = WriterStatusDone
	return 1, nil

}

func (w *Writer) WriteTrailers(headers headers.Header) error {
	if w.WriteStatus != WriterStatusWritingBody {
		return fmt.Errorf("cannot write trailers in current state: %v", w.WriteStatus)
	}
	if err := WriteHeaders(w, headers); err != nil {
		w.WriteStatus = WriterStatusError
		return fmt.Errorf("error writing trailers: %w", err)
	}
	return nil
}

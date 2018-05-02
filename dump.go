package gindump

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	dumpWriter io.Writer
	gin.ResponseWriter
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) WriteHeaderNow() {
	w.ResponseWriter.WriteHeaderNow()
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.dumpWriter.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *responseWriter) Size() int {
	return w.ResponseWriter.Size()
}

func (w *responseWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.Hijack()
}

func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.CloseNotify()
}

func (w *responseWriter) Flush() {
	w.ResponseWriter.Flush()
}

const (
	keyRequest  = "key_dumped_request"
	keyResponse = "key_dumped_response"
)

func Dump() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		reqBody := new(bytes.Buffer)
		r := io.TeeReader(c.Request.Body, reqBody)
		c.Request.Body = ioutil.NopCloser(r)

		// response
		resBody := new(bytes.Buffer)
		w := &responseWriter{
			dumpWriter:     resBody,
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		c.Next()

		req := &http.Request{
			Method: c.Request.Method,
			URL:    c.Request.URL,
			Proto:  c.Request.Proto,
			Header: c.Request.Header,
			Body:   ioutil.NopCloser(bytes.NewBuffer(reqBody.Bytes())),
		}

		res := &http.Response{
			StatusCode: c.Writer.Status(),
			Header:     c.Writer.Header(),
			Body:       ioutil.NopCloser(bytes.NewBuffer(resBody.Bytes())),
		}

		c.Set(keyRequest, req)
		c.Set(keyResponse, res)
	}
}

type DumpRequest struct {
	Method string
	URL    *url.URL
	Proto  string
	Header http.Header
	Body   io.ReadCloser
}

type DumpResponse struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

func GetRequest(c *gin.Context) (*DumpRequest, error) {
	v, ok := c.Get(keyRequest)
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	r, ok := v.(*http.Request)
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	return &DumpRequest{
		Method: r.Method,
		URL:    r.URL,
		Proto:  r.Proto,
		Header: r.Header,
		Body:   r.Body,
	}, nil
}

func GetResponse(c *gin.Context) (*DumpResponse, error) {
	v, ok := c.Get(keyResponse)
	if !ok {
		return nil, fmt.Errorf("response not found")
	}
	r, ok := v.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("response not found")
	}
	return &DumpResponse{
		StatusCode: r.StatusCode,
		Header:     r.Header,
		Body:       r.Body,
	}, nil
}

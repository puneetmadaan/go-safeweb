// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package safehttp

import (
	"net/http"
)

// ResponseWriter TODO
type ResponseWriter struct {
	d  Dispatcher
	rw http.ResponseWriter

	// Having this field unexported is essential for
	// security. Otherwise one can easily overwrite
	// the struct bypassing all our safety guarantees.
	header Header
}

// NewResponseWriter creates a ResponseWriter from
// a safehttp.Dispatcher and a http.ResponseWriter.
func NewResponseWriter(d Dispatcher, rw http.ResponseWriter) ResponseWriter {
	header := newHeader(rw.Header())
	return ResponseWriter{d: d, rw: rw, header: header}
}

// Result TODO
type Result struct{}

// Write TODO
func (w *ResponseWriter) Write(resp Response) Result {
	if err := w.d.Write(w.rw, resp); err != nil {
		panic("error")
	}
	return Result{}
}

// WriteTemplate TODO
func (w *ResponseWriter) WriteTemplate(t Template, data interface{}) Result {
	if err := w.d.ExecuteTemplate(w.rw, t, data); err != nil {
		panic("error")
	}
	return Result{}
}

// ClientError TODO
func (w *ResponseWriter) ClientError(code StatusCode) Result {
	if code < 400 || code >= 500 {
		// TODO(@mihalimara22): Replace this when we decide how to handle this case
		panic("wrong method called")
	}
	http.Error(w.rw, http.StatusText(int(code)), int(code))
	return Result{}
}

// ServerError TODO
func (w *ResponseWriter) ServerError(code StatusCode) Result {
	if code < 500 || code >= 600 {
		// TODO(@mattiasgrenfeldt, @mihalimara22, @kele, @empijei): Decide how it should
		// be communicated to the user of the framework that they've called the wrong
		// method.
		return Result{}
	}
	http.Error(w.rw, http.StatusText(int(code)), int(code))
	return Result{}
}

// Redirect responds with a redirect to a given url, using code as the status code.
func (w *ResponseWriter) Redirect(r *IncomingRequest, url string, code StatusCode) Result {
	http.Redirect(w.rw, r.req, url, int(code))
	return Result{}
}

// Header returns the collection of headers that will be set
// on the response. Headers must be set before writing a
// response (e.g. Write, WriteTemplate).
func (w ResponseWriter) Header() Header {
	return w.header
}

// Dispatcher TODO
type Dispatcher interface {
	Write(rw http.ResponseWriter, resp Response) error
	ExecuteTemplate(rw http.ResponseWriter, t Template, data interface{}) error
}

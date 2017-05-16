package lib

import (
	"io"
	"net/http"
)

// Finalizer finalizes an htpp response.
type Finalizer interface {
	HandleError(err error, w io.Writer, r *http.Request) bool
	HandleSuccess(w io.Writer, r io.Reader) error
}

// DefaultFinalizer for an http response.
type DefaultFinalizer struct {
}

// HandleError prints http 500.
func (f DefaultFinalizer) HandleError(err error, w io.Writer, r *http.Request) bool {
	if x, ok := w.(http.ResponseWriter); ok {
		http.Error(x, err.Error(), http.StatusInternalServerError)
	}
	return true
}

// HandleSuccess prints http 200 and prints r.
func (f DefaultFinalizer) HandleSuccess(w io.Writer, r io.Reader) error {
	if x, ok := w.(http.ResponseWriter); ok {
		x.WriteHeader(http.StatusOK)
	}
	return nil
}

// JSONFinalizer finalizes a JSON response.
type JSONFinalizer struct {
	DefaultFinalizer
}

// HandleSuccess prints http 200 and prints r.
func (f JSONFinalizer) HandleSuccess(w io.Writer, r io.Reader) error {
	f.DefaultFinalizer.HandleSuccess(w, r)
	if x, ok := w.(http.ResponseWriter); ok {
		x.Header().Set("Content-Type", "application/json")
	}
	_, err := io.Copy(w, r)
	return err
}

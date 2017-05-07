package demo

// file generated by
// github.com/mh-cbon/jsoner
// do not edit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// JSONTomates is jsoner of *Tomates.
type JSONTomates struct {
	embed *Tomates
}

// NewJSONTomates constructs a jsoner of *Tomates
func NewJSONTomates(embed *Tomates) *JSONTomates {
	ret := &JSONTomates{
		embed: embed,
	}
	return ret
}

// HandleSuccess prints http 200 and prints r.
func (t *JSONTomates) HandleSuccess(w io.Writer, r io.Reader) error {
	if x, ok := w.(http.ResponseWriter); ok {
		x.WriteHeader(http.StatusOK)
		x.Header().Set("Content-Type", "application/json")
	}
	_, err := io.Copy(w, r)
	return err
}

// Push reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Push(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		x []Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar0 := t.embed.Push(input.x...)

	out, encErr := json.Marshal([]interface{}{retVar0})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Unshift reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Unshift(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		x []Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar1 := t.embed.Unshift(input.x...)

	out, encErr := json.Marshal([]interface{}{retVar1})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Pop reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Pop(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar2 := t.embed.Pop()

	out, encErr := json.Marshal([]interface{}{retVar2})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Shift reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Shift(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar3 := t.embed.Shift()

	out, encErr := json.Marshal([]interface{}{retVar3})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Index reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Index(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		s Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar4 := t.embed.Index(input.s)

	out, encErr := json.Marshal([]interface{}{retVar4})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Contains reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Contains(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		s Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar5 := t.embed.Contains(input.s)

	out, encErr := json.Marshal([]interface{}{retVar5})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// RemoveAt reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) RemoveAt(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		i int
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar6 := t.embed.RemoveAt(input.i)

	out, encErr := json.Marshal([]interface{}{retVar6})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Remove reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Remove(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		s Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar7 := t.embed.Remove(input.s)

	out, encErr := json.Marshal([]interface{}{retVar7})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// InsertAt reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) InsertAt(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		i int
		s Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar8 := t.embed.InsertAt(input.i, input.s)

	out, encErr := json.Marshal([]interface{}{retVar8})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Splice reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Splice(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		start  int
		length int
		s      []Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar9 := t.embed.Splice(input.start, input.length, input.s...)

	out, encErr := json.Marshal([]interface{}{retVar9})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Slice reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Slice(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		start  int
		length int
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar10 := t.embed.Slice(input.start, input.length)

	out, encErr := json.Marshal([]interface{}{retVar10})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Reverse reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Reverse(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar11 := t.embed.Reverse()

	out, encErr := json.Marshal([]interface{}{retVar11})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Len reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Len(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar12 := t.embed.Len()

	out, encErr := json.Marshal([]interface{}{retVar12})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Set reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Set(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		x []Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar13 := t.embed.Set(input.x)

	out, encErr := json.Marshal([]interface{}{retVar13})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// Get reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) Get(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar14 := t.embed.Get()

	out, encErr := json.Marshal([]interface{}{retVar14})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

// At reads json, outputs json.
// the json input must provide a key/value for each params.
func (t *JSONTomates) At(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	input := struct {
		i int
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar15 := t.embed.At(input.i)

	out, encErr := json.Marshal([]interface{}{retVar15})
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(out)
		ret = &b
	}

	return ret, retErr

}

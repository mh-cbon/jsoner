package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//go:generate lister Tomate:TomatesGen
//go:generate jsoner *TomatesGen:TomatesJSONGen
//go:generate jsoner *Controller:ControllerJSONGen

func main() {
	// setup json rpc calls on a slice of tomates
	tomatesSlice := NewTomatesGen()
	tomatesJSON := NewTomatesJSONGen(tomatesSlice, nil)

	// make a push req on the slice.
	r := makeJSONRPCPushReq(Tomate{Name: "Hello world!"})

	// apply the req on the json rpcer to the slice of tomates.
	ret, err := tomatesJSON.Push(r)
	if err != nil {
		panic(err)
	}

	// decode the json response
	res := struct{ Arg0 *TomatesGen }{}
	json.Unmarshal(ret.(*bytes.Buffer).Bytes(), &res)
	fmt.Println(res.Arg0.First().Name)
}

func makeJSONRPCPushReq(t Tomate) *http.Request {
	r := &http.Request{}

	params := struct {
		Arg0 []Tomate
	}{
		Arg0: []Tomate{t},
	}

	b, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	var buf closeBuffer
	(&buf).Write(b)
	r.Body = &buf

	return r
}

type closeBuffer struct{ bytes.Buffer } // not sure why i needed to do that, there s must be better ways.

func (c closeBuffer) Close() error { return nil }

// Tomate if the resource subject.
type Tomate struct {
	Name string
}

// GetID ...
func (t Tomate) GetID() string {
	return t.Name
}

// Controller of some resources.
type Controller struct {
}

// GetByID ...
func (t Controller) GetByID(id int) Tomate {
	return Tomate{}
}

// UpdateByID ...
func (t Controller) UpdateByID(GETid int, reqBody Tomate) Tomate {
	return Tomate{}
}

// DeleteByID ...
func (t *Controller) DeleteByID(reqID int) bool {
	return false
}

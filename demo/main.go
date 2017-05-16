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

	//====================================================
	//                                              demo 1

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

	//====================================================
	//                                              demo 2

	// setup json rpc calls on a tomate slice controller
	tomateCtl := NewController(tomatesSlice)
	JSONTomateCtl := NewControllerJSONGen(tomateCtl, nil)

	// add new value for the demo
	tomatesSlice.Push(Tomate{"Red"})

	// make a GetByName req on the slice.
	r2 := makeJSONRPCGetByNameReq("Red")

	// apply the req on the json rpcer to the slice of tomates.
	ret2, err2 := JSONTomateCtl.GetByName(r2)
	if err2 != nil {
		panic(err2)
	}

	// decode the json response
	res2 := struct{ Arg0 Tomate }{}
	json.Unmarshal(ret2.(*bytes.Buffer).Bytes(), &res2)
	fmt.Println(res2.Arg0.Name)
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

func makeJSONRPCGetByNameReq(n string) *http.Request {
	r := &http.Request{}

	params := struct {
		Arg0 string
	}{
		Arg0: n,
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

// NewController is a constructor.
func NewController(backend *TomatesGen) *Controller {
	return &Controller{backend: backend}
}

// Controller of some resources.
type Controller struct {
	backend *TomatesGen
}

// GetByName ...
func (t *Controller) GetByName(name string) Tomate {
	return t.backend.Filter(FilterTomatesGen.ByName(name)).First()
}

// UpdateByName ...
func (t *Controller) UpdateByName(GETname string, reqBody Tomate) Tomate {
	t.backend.Map(func(x Tomate) Tomate {
		if x.Name == GETname {
			return reqBody
		}
		return x
	})
	return reqBody
}

// DeleteByName ...
func (t *Controller) DeleteByName(reqName string) bool {
	return t.backend.Remove(Tomate{Name: reqName})
}

# jsoner

[![travis Status](https://travis-ci.org/mh-cbon/jsoner.svg?branch=master)](https://travis-ci.org/mh-cbon/jsoner) [![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/jsoner?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/jsoner) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/jsoner)](https://goreportcard.com/report/github.com/mh-cbon/jsoner) [![GoDoc](https://godoc.org/github.com/mh-cbon/jsoner?status.svg)](http://godoc.org/github.com/mh-cbon/jsoner) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Package jsoner is a cli tool to implement json-rpc of a type.


s/Choose your gun!/[Aux armes!](https://www.youtube.com/watch?v=hD-wD_AMRYc&t=7)/

# TOC
- [Install](#install)
  - [Usage](#usage)
    - [$ jsoner -help](#-jsoner--help)
  - [Cli examples](#cli-examples)
- [API example](#api-example)
  - [> demo/main.go](#-demomaingo)
  - [> demo/controllerjsongen.go](#-democontrollerjsongengo)
  - [Conventionned variable name](#conventionned-variable-name)
- [Recipes](#recipes)
  - [Release the project](#release-the-project)
- [History](#history)

# Install
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/jsoner
cd $GOPATH/src/github.com/mh-cbon/jsoner
git clone https://github.com/mh-cbon/jsoner.git .
glide install
go install
```

## Usage

#### $ jsoner -help
```sh
jsoner 0.0.0

Usage

  jsoner [-p name] [...types]

  types:  A list of types such as src:dst.
          A type is defined by its package path and its type name,
          [pkgpath/]name
          If the Package path is empty, it is set to the package name being generated.
          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName 
  -p:     The name of the package output.
```

## Cli examples

```sh
# Create a jsoned version of Tomate to MyTomate
jsoner tomate_gen.go Tomate:MyTomate
# Create a jsoned version of Tomate to MyTomate to stdout
lister -p main - Tomate:MyTomate
```

# API example

Following example demonstates a program using it to generate a jsoned version of a type.

#### > demo/main.go
```go
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
	// 																							demo 1

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
	// 																							demo 2

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
```

Following code is the generated implementation of a jsoner `Controller` that uses conventionend variable name.

#### > demo/controllerjsongen.go
```go
package main

// file generated by
// github.com/mh-cbon/jsoner
// do not edit

import (
	"bytes"
	"encoding/json"
	jsoner "github.com/mh-cbon/jsoner/lib"
	"io"
	"net/http"
)

// ControllerJSONGen is jsoner of *Controller.
// Controller of some resources.
type ControllerJSONGen struct {
	embed     *Controller
	finalizer jsoner.Finalizer
}

// NewControllerJSONGen constructs a jsoner of *Controller
func NewControllerJSONGen(embed *Controller, finalizer jsoner.Finalizer) *ControllerJSONGen {
	if finalizer == nil {
		finalizer = &jsoner.JSONFinalizer{}
	}
	ret := &ControllerJSONGen{
		embed:     embed,
		finalizer: finalizer,
	}
	return ret
}

//UnmarshalJSON JSON unserializes ControllerJSONGen
func (t *ControllerJSONGen) UnmarshalJSON(b []byte) error {
	var embed *Controller
	if err := json.Unmarshal(b, &embed); err != nil {
		return err
	}
	t.embed = embed
	return nil
}

//MarshalJSON JSON serializes ControllerJSONGen
func (t *ControllerJSONGen) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.embed)
}

// GetByName Decodes r as json to invoke *Controller.GetByName.
// GetByName ...
func (t *ControllerJSONGen) GetByName(r *http.Request) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error
	input := struct {
		Arg0 string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)
	if decErr != nil {
		return nil, decErr
	}

	retVar0 := t.embed.GetByName(input.Arg0)

	output := struct {
		Arg0 Tomate
	}{
		Arg0: retVar0,
	}

	outBytes, encErr := json.Marshal(output)
	if encErr != nil {
		retErr = encErr
	} else {
		var b bytes.Buffer
		b.Write(outBytes)
		ret = &b
	}

	return ret, retErr

}

// UpdateByName Decodes reqBody as json to invoke *Controller.UpdateByName.
// Other parameters are passed straight
// UpdateByName ...
func (t *ControllerJSONGen) UpdateByName(GETname string, reqBody io.Reader) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	var decBody Tomate
	decErr := json.NewDecoder(reqBody).Decode(&decBody)
	if decErr != nil {
		return nil, decErr
	}
	retVar1 := t.embed.UpdateByName(GETname, decBody)

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

// DeleteByName Decodes reqBody as json to invoke *Controller.DeleteByName.
// Other parameters are passed straight
// DeleteByName ...
func (t *ControllerJSONGen) DeleteByName(reqName string) (io.Reader, error) {

	ret := new(bytes.Buffer)
	var retErr error

	retVar2 := t.embed.DeleteByName(reqName)

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
```

#### Conventionned variable name

`jsoner` reads and interprets input params to determine where/how the request body should be decoded.

In that matters it looks up for variable named `reqBody`,
if such variable is found, then the body is decoded according to its type.
Other parameters are passed/received untouched.

If such variable is not found, it is assumed that all parameters are to be decoded from the request body.

# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)

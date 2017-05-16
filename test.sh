#!/bin/sh

set -ex

rm `which jsoner`
go install


jsoner - demo/Controller:ControllerJSON | grep -F "GetByID(r *http.Request" || exit 1;
jsoner - demo/Controller:ControllerJSON | grep "package main" || exit 1;
jsoner -p nop - demo/Controller:ControllerJSON | grep "package nop" || exit 1;

jsoner - demo/Controller:ControllerJSON | grep "embed Controller" || exit 1;
jsoner - demo/*Controller:ControllerJSON | grep -F "embed *Controller" || exit 1;

jsoner - demo/Controller:*ControllerJSON | grep "embed Controller" || exit 1;
jsoner - demo/*Controller:*ControllerJSON | grep -F "embed *Controller" || exit 1;

rm -fr gen_test
jsoner demo/Controller:gen_test/ControllerJSON || exit 1;
ls -al gen_test | grep "controllerjson.go" || exit 1;
cat gen_test/controllerjson.go | grep -F "GetByID(r *http.Request" || exit 1;
cat gen_test/controllerjson.go | grep "package gen_test" || exit 1;
rm -fr gen_test

rm -fr demo/*gen.go
go generate demo/main.go
ls -al demo | grep "controllerjsongen.go" || exit 1;
cat demo/controllerjsongen.go | grep "package main" || exit 1;
cat demo/controllerjsongen.go | grep "NewControllerJSONGen(" || exit 1;
cat demo/tomatesjsongen.go | grep "NewTomatesJSONGen(" || exit 1;
go run demo/*.go | grep "Hello world!" || exit 1;

rm -fr demo/*gen.go
go generate github.com/mh-cbon/jsoner/demo
ls -al demo | grep "tomatesjsongen.go" || exit 1;
ls -al demo | grep "controllerjsongen.go" || exit 1;
cat demo/controllerjsongen.go | grep "package main" || exit 1;
cat demo/controllerjsongen.go | grep "NewControllerJSONGen(" || exit 1;
cat demo/tomatesjsongen.go | grep "NewTomatesJSONGen(" || exit 1;
go run demo/*.go | grep "Hello world!" || exit 1;
# # rm -fr demo/gen # keep it for demo

# go test


echo ""
echo "ALL GOOD!"

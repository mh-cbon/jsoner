---
License: MIT
LicenseFile: LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "badge/travis" .}} {{template "badge/appveyor" .}} {{template "badge/goreport" .}} {{template "badge/godoc" .}} {{template "license/shields" .}}

{{pkgdoc}}

Choose your gun! | [Aux armes!](https://www.youtube.com/watch?v=hD-wD_AMRYc&t=7)

# {{toc 5}}

# Install
{{template "glide/install" .}}

## Usage

#### $ {{exec "jsoner" "-help" | color "sh"}}

## Cli examples

```sh
# Create a jsoned version of Tomate to MyTomate
jsoner tomate_gen.go Tomate:MyTomate
```
# API example

Following example demonstates a program using it to generate a jsoned version of a type.

#### > {{cat "demo/lib.go" | color "go"}}

Following code is the generated implementation of a typed slice of `Tomate`.

#### > {{cat "demo/json_vegetables_gen.go" | color "go"}}

Following code is the generated implementation of a jsoner `Controller` that uses conventionend variable name.

#### > {{cat "demo/json_controller_gen.go" | color "go"}}

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

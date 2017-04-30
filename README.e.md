---
License: MIT
LicenseFile: LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "badge/travis" .}} {{template "badge/appveyor" .}} {{template "badge/goreport" .}} {{template "badge/godoc" .}} {{template "license/shields" .}}

{{pkgdoc}}

Get your weapon! https://www.youtube.com/watch?v=hD-wD_AMRYc&t=7

# {{toc 5}}

# Install
{{template "glide/install" .}}

## Usage

#### $ {{exec "jsoner" "-help" | color "sh"}}

## Cli examples

```sh
# Create a channeled version of Tomate to MyTomate
jsoner tomate_gen.go Tomate:MyTomate
```
# API example

Following example demonstates a program using it to generate a channeled version of a type.

#### > {{cat "demo/lib.go" | color "go"}}

Following code is the generated implementation of a typed slice of `Tomate`.

#### > {{cat "demo/json_vegetables_gen.go" | color "go"}}

# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)

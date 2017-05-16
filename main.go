// Package jsoner is a cli tool to implement json-rpc of a type.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mh-cbon/astutil"
	httper "github.com/mh-cbon/httper/lib"
	"github.com/mh-cbon/jsoner/utils"
)

var name = "jsoner"
var version = "0.0.0"

func main() {

	var help bool
	var h bool
	var ver bool
	var v bool
	var outPkg string
	var mode string
	flag.BoolVar(&help, "help", false, "Show help.")
	flag.BoolVar(&h, "h", false, "Show help.")
	flag.BoolVar(&ver, "version", false, "Show version.")
	flag.BoolVar(&v, "v", false, "Show version.")
	flag.StringVar(&outPkg, "p", os.Getenv("GOPACKAGE"), "Package name of the new code.")
	flag.StringVar(&mode, "mode", "std", "Generation mode.")

	flag.Parse()

	if ver || v {
		showVer()
		return
	}
	if help || h {
		showHelp()
		return
	}

	if flag.NArg() < 1 {
		panic("wrong usage")
	}
	args := flag.Args()

	out := ""
	if args[0] == "-" {
		args = args[1:]
		out = "-"
	}

	todos, err := utils.NewTransformsArgs(utils.GetPkgToLoad()).Parse(args)
	if err != nil {
		panic(err)
	}

	filesOut := utils.NewFilesOut("github.com/mh-cbon/" + name)

	for _, todo := range todos.Args {
		if todo.FromPkgPath == "" {
			log.Println("Skipped ", todo.FromTypeName)
			continue
		}

		fileOut := filesOut.Get(todo.ToPath)

		fileOut.PkgName = outPkg
		if fileOut.PkgName == "" {
			fileOut.PkgName = findOutPkg(todo)
		}

		if err := processType(mode, todo, fileOut); err != nil {
			log.Println(err)
		}
	}

	filesOut.Write(out)
}

func showVer() {
	fmt.Printf("%v %v\n", name, version)
}

func showHelp() {
	showVer()
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println()
	fmt.Printf("  %v [-p name] [-mode name] [...types]\n\n", name)
	fmt.Printf("  types:  A list of types such as src:dst.\n")
	fmt.Printf("          A type is defined by its package path and its type name,\n")
	fmt.Printf("          [pkgpath/]name\n")
	fmt.Printf("          If the Package path is empty, it is set to the package name being generated.\n")
	// fmt.Printf("          If the Package path is a directory relative to the cwd, and the Package name is not provided\n")
	// fmt.Printf("          the package path is set to this relative directory,\n")
	// fmt.Printf("          the package name is set to the name of this directory.\n")
	fmt.Printf("          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName \n")
	fmt.Printf("  -p:     The name of the package output.\n")
	fmt.Println()
}

func findOutPkg(todo utils.TransformArg) string {
	if todo.ToPkgPath != "" {
		prog := astutil.GetProgramFast(todo.ToPkgPath)
		if prog != nil {
			pkg := prog.Package(todo.ToPkgPath)
			return pkg.Pkg.Name()
		}
	}
	if todo.ToPkgPath == "" {
		prog := astutil.GetProgramFast(utils.GetPkgToLoad())
		if len(prog.Imported) < 1 {
			panic("impossible, add [-p name] option")
		}
		for _, p := range prog.Imported {
			return p.Pkg.Name()
		}
	}
	if strings.Index(todo.ToPkgPath, "/") > -1 {
		return filepath.Base(todo.ToPkgPath)
	}
	return todo.ToPkgPath
}

func processType(mode string, todo utils.TransformArg, fileOut *utils.FileOut) error {

	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)
	foundMethods := astutil.FindMethods(pkg)

	if todo.FromPkgPath != todo.ToPkgPath {
		fileOut.AddImport(todo.FromPkgPath, "")
	}
	if todo.FromPkgPath != todo.ToPkgPath {
		fileOut.AddImport(todo.FromPkgPath, "")
	}
	fileOut.AddImport("bytes", "")
	fileOut.AddImport("encoding/json", "")
	fileOut.AddImport("io", "")
	fileOut.AddImport("net/http", "")
	fileOut.AddImport("github.com/mh-cbon/jsoner/lib", "jsoner")

	// func processType(mode, destName, srcName string, prog *loader.Program, pkg *loader.PackageInfo, foundMethods map[string][]*ast.FuncDecl) ([]string, bytes.Buffer) {

	// extraImports := []string{}
	// var b bytes.Buffer
	// dest := &b

	srcConcrete := astutil.GetUnpointedType(srcName)
	dstConcrete := astutil.GetUnpointedType(destName)

	structType := astutil.FindStruct(pkg, srcConcrete)

	if structType == nil {
		panic(
			fmt.Errorf("Type not found: %v", srcName),
		)
	}

	structComment := astutil.GetComment(prog, structType.Pos())
	// todo: might do better to send only annotations or do other improvemenets.
	structComment = makeCommentLines(structComment)

	fmt.Fprintf(dest, `
// %v is jsoner of %v.
%v
type %v struct{
	embed %v
	finalizer jsoner.Finalizer
}
		`, dstConcrete, srcName, structComment, dstConcrete, srcName)

	dstStar := astutil.GetPointedType(destName)
	// hasHandleError := methodsContains(srcName, "HandleError", foundMethods)
	// hasHandleSuccess := methodsContains(srcName, "HandleSuccess", foundMethods)

	// Make the constructor
	fmt.Fprintf(dest, `// New%v constructs a jsoner of %v
func New%v(embed %v, finalizer jsoner.Finalizer) %v {
	if finalizer == nil {
		finalizer = &jsoner.JSONFinalizer{}
	}
	ret := &%v{
		embed: embed,
		finalizer: finalizer,
	}
  return ret
}
`, dstConcrete, srcName, dstConcrete, srcName, dstStar, dstConcrete)
	fmt.Fprintln(dest)

	// Add marshalling capabilities
	fmt.Fprintf(dest, `
//UnmarshalJSON JSON unserializes %v
func (t %v) UnmarshalJSON(b []byte) error {
	var embed %v
	if err := json.Unmarshal(b, &embed); err != nil {
		return err
	}
	t.embed = embed
	return nil
}
`, dstConcrete, dstStar, srcName)
	fmt.Fprintln(dest)

	fmt.Fprintf(dest, `
//MarshalJSON JSON serializes %v
func (t %v) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.embed)
}
`, dstConcrete, dstStar)
	fmt.Fprintln(dest)

	for _, m := range foundMethods[srcConcrete] {

		methodName := astutil.MethodName(m)
		params := astutil.MethodParams(m)
		paramNames := astutil.MethodParamNames(m)
		retTypes := astutil.MethodReturnTypes(m)
		retVars := astutil.MethodReturnVars(m)
		sRetVars := strings.Join(retVars, ", ")
		sRetVars = strings.TrimSpace(sRetVars)
		hasErr := astutil.MethodReturnError(m)
		// structProps := astutil.MethodParamsToProps(m)

		importIDs := astutil.GetSignatureImportIdentifiers(m)
		for _, i := range importIDs {
			fileOut.AddImport(i, "")
		}
		// receiverName := astutil.ReceiverName(m)
		comment := astutil.GetComment(prog, m.Pos())
		comment = makeCommentLines(comment)
		lParams := strings.Split(params, ",")
		lParamNames := strings.Split(paramNames, ",")
		hasEllipse := astutil.MethodHasEllipse(m)

		// ensure it is desired to facade this method.
		if astutil.IsExported(methodName) == false {
			continue
		}

		// verify that the method does not take unserializable arguments.
		// todo: improve marshaller support detection.
		if !isMarshable(params) {
			break
		}

		if !isUsingConvetionnedParams(mode, params) {
			// the method params does not use
			// conventionned params names,
			// the parameters must be decoded from the
			// req body, and applied to the method.
			methInvok := fmt.Sprintf(`t.embed.%v()
					`, methodName)

			if sRetVars != "" {
				methInvok = sRetVars + ":=" + methInvok
			}

			if params != "" {

				methInvok = fmt.Sprintf(`input := struct{
						%v
					}{}
					decErr := json.NewDecoder(r.Body).Decode(&input)
					if decErr != nil {
						return nil, decErr
					}
				`, mapParamsToStruct(lParams, hasEllipse))

				methInvok += fmt.Sprintf(`
					%v := t.embed.%v(%v)
					`, sRetVars, methodName, mapParamNamesToStructProps(lParamNames, hasEllipse))
			}

			errHandling := ""
			if hasErr {
				errName := retVars[len(retVars)-1]
				errHandling = fmt.Sprintf(`if %v != nil {
						retErr = %v
					}`, errName, errName)
			}

			outHandling := ""
			if sRetVars != "" {
				outHandling = fmt.Sprintf(`
					output := struct{
						%v
					}{
						%v
					}
					`, mapParamsToStruct(retTypes, false), mapParamsToStructValues(retVars))

				outHandling += fmt.Sprint(`
					outBytes, encErr := json.Marshal(output)
					if encErr!= nil {
						retErr = encErr
					} else {
						var b bytes.Buffer
						b.Write(outBytes)
						ret = &b
					}
					`)
			}

			body := fmt.Sprintf(`
				ret := new(bytes.Buffer)
				var retErr error
				%v
				%v
				%v
				return ret, retErr
				`, methInvok, errHandling, outHandling)

			fmt.Fprintf(dest, `// %v Decodes r as json to invoke %v.%v.
				%v
				`, methodName, srcName, methodName, comment)
			fmt.Fprintf(dest, `func (t %v) %v(r *http.Request) (io.Reader, error) {
					%v
				}
				`, dstStar, methodName, body)

		} else {
			// the method params uses
			// conventionned params names,
			// the req body should be decoded according to reqBody param,
			// other parameters are to be received/forwarded regularly.

			// the reqBody param type is set to some concrete type by the end user,
			// jsoner should decode the request body to that concrete type.

			// set the type of the reqBody of the new method to io.Reader,
			// this is what a protocol implementation d probably provide.
			targetType := getParamType(lParams, "reqBody")
			newlParams := changeParamType(lParams, "reqBody", "io.Reader")
			newParams := strings.Join(newlParams, ",")

			// build a new list of param invokation where the input reqBody is replaced by decBody
			newlParamNames := changeParamName(lParamNames, "reqBody", "decBody")
			newParamNames := strings.Join(newlParamNames, ",")

			// decode the reqBody into decBody
			bodyDec := ""
			if targetType != "" {
				amp := "" // & ?
				if !astutil.IsAPointedType(targetType) &&
					!astutil.IsASlicedType(targetType) {
					amp = "&"
				}
				bodyDec = fmt.Sprintf(`
					var decBody %v
					decErr := json.NewDecoder(reqBody).Decode(%vdecBody)
					if decErr != nil {
						return nil, decErr
					}`, targetType, amp)
			}

			// invoke the embeded method with the new params list.
			methInvok := fmt.Sprintf(`t.embed.%v(%v)
				`, methodName, newParamNames)

			if strings.TrimSpace(sRetVars) != "" {
				methInvok = sRetVars + ":=" + methInvok
			}

			errHandling := ""
			if hasErr {
				errName := retVars[len(retVars)-1]
				errHandling = fmt.Sprintf(`if %v != nil {
								retErr = %v
				}`, errName, errName)
			}

			outHandling := ""
			if sRetVars != "" {
				outHandling = fmt.Sprintf(`
					out, encErr := json.Marshal([]interface{}{%v})
					if encErr!= nil {
						retErr = encErr
					} else {
						var b bytes.Buffer
						b.Write(out)
						ret = &b
					}
						`, sRetVars)
			}

			body := fmt.Sprintf(`
				ret := new(bytes.Buffer)
				var retErr error
				%v
				%v
				%v
				%v
				return ret, retErr
				`, bodyDec, methInvok, errHandling, outHandling)

			fmt.Fprintf(dest, `// %v Decodes reqBody as json to invoke %v.%v.
				// Other parameters are passed straight
				%v
				`, methodName, srcName, methodName, comment)

			fmt.Fprintf(dest, `func (t %v) %v(%v) (io.Reader, error) {
					%v
				}
				`, dstStar, methodName, newParams, body)
		}
	}

	return nil
}

func getParamType(lParams []string, name string) string {
	for _, p := range lParams {
		p = strings.TrimSpace(p)
		if strings.Index(p, name) == 0 {
			return strings.Split(p, " ")[1]
		}
	}
	return ""
}

func changeParamType(lParams []string, name, t string) []string {
	ret := []string{}
	for _, p := range lParams {
		p = strings.TrimSpace(p)
		if strings.Index(p, name) == 0 {
			p = name + " " + t
		}
		ret = append(ret, p)
	}
	return ret
}

func changeParamName(lParamNames []string, name, to string) []string {
	ret := []string{}
	for _, p := range lParamNames {
		p = strings.TrimSpace(p)
		if strings.Index(p, name) == 0 {
			p = to + " " + p[len(name):]
		}
		ret = append(ret, p)
	}
	return ret
}

func mapParamsToStruct(params []string, hasEllipse bool) string {
	ret := ""
	if len(params) > 0 {
		for i, p := range params {
			p = strings.TrimSpace(p)
			y := strings.Split(p, " ")
			t := strings.Replace(y[0], "...", "", -1)
			if len(y) > 1 {
				t = strings.Replace(y[1], "...", "", -1)
			}
			if i == len(params)-1 && hasEllipse {
				ret += fmt.Sprintf("Arg%v []%v\n", i, t)
			} else {
				ret += fmt.Sprintf("Arg%v %v\n", i, t)
			}
		}
	}
	return ret
}

func mapParamsToStructValues(params []string) string {
	ret := ""
	if len(params) > 0 {
		for i, p := range params {
			p = strings.TrimSpace(p)
			y := strings.Split(p, " ")
			ret += fmt.Sprintf("Arg%v: %v,\n", i, y[0])
		}
	}
	return ret
}

func mapParamNamesToStructProps(params []string, hasEllipse bool) string {
	ret := ""
	if len(params) > 0 {
		for i := range params {
			if i == len(params)-1 && hasEllipse {
				ret += fmt.Sprintf("input.Arg%v..., ", i)
			} else {
				ret += fmt.Sprintf("input.Arg%v, ", i)
			}
		}
		ret = ret[:len(ret)-2]
	}
	return ret
}

func makeCommentLines(s string) string {
	s = strings.TrimSpace(s)
	comment := ""
	for _, k := range strings.Split(s, "\n") {
		comment += "// " + k + "\n"
	}
	comment = strings.TrimSpace(comment)
	if comment == "" {
		comment = "//"
	}
	return comment
}

var gorillaMode = "gorilla"
var stdMode = "std"
var reqBodyVarName = "reqBody"

func isUsingConvetionnedParams(mode, params string) bool {
	lParams := strings.Split(params, ",")
	for _, param := range lParams {
		k := strings.Split(param, " ")
		if len(k) > 1 {
			varType := strings.TrimSpace(k[1])
			if varType == "http.ResponseWriter" {
				return true

			} else if varType == "*http.Request" {
				return true

			} else if varType == "httper.Cookier" {
				return true

			} else if varType == "httper.Sessionner" {
				return true
			}
		}
		varName := strings.TrimSpace(k[0])
		if isConvetionnedParam(mode, varName) {
			return true
		}
	}
	return false
}

func isConvetionnedParam(mode, varName string) bool {
	if varName == reqBodyVarName {
		return true
	}
	return getVarPrefix(mode, varName) != ""
}

func getParamConvention(mode, varName string) string {
	if varName == reqBodyVarName {
		return reqBodyVarName
	}
	return getVarPrefix(mode, varName)
}

func getSessionProviderFactory(mode string) httper.SessionProvider {
	var factory httper.SessionProvider
	if mode == stdMode {
		factory = &httper.VoidSessionProvider{}
	} else if mode == gorillaMode {
		factory = &httper.GorillaSessionProvider{}
	}
	return factory
}

func getDataProviderFactory(mode string) httper.DataerProvider {
	var factory httper.DataerProvider
	if mode == stdMode {
		factory = &httper.StdHTTPDataProvider{}
	} else if mode == gorillaMode {
		factory = &httper.GorillaHTTPDataProvider{}
	}
	return factory
}

func getDataProvider(mode string) *httper.DataProviderFacade {
	return getDataProviderFactory(mode).MakeEmpty().(*httper.DataProviderFacade)
}

func getVarPrefix(mode, varName string) string {
	ret := ""
	provider := getDataProvider(mode)
	for _, p := range provider.Providers {
		prefix := p.GetName()
		if strings.HasPrefix(varName, strings.ToLower(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToUpper(f) {
				ret = prefix
				break
			}
		} else if strings.HasPrefix(varName, strings.ToUpper(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToLower(f) {
				ret = prefix
				break
			}
		}
	}
	return ret
}

func methodsContains(typeName, search string, methods map[string][]*ast.FuncDecl) bool {
	if funList, ok := methods[typeName]; ok {
		for _, fun := range funList {
			if astutil.MethodName(fun) == search {
				return true
			}
		}
	}
	return false
}

func isMarshable(params string) bool {
	ret := true
	for _, p := range strings.Split(params, ",") {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			p = strings.Join(strings.Split(p, " ")[1:], " ") // get the type
			p = strings.TrimSpace(p)
			if strings.Index(p, "func") > -1 {
				ret = false // if its a param of type func
				break
			} else if strings.Index(p, "chan") > -1 {
				ret = false // if its a param of type chan
				break
			}
		}
	}
	return ret
}

func getPkgToLoad() string {
	gopath := filepath.Join(os.Getenv("GOPATH"), "src")
	pkgToLoad, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pkgToLoad[len(gopath)+1:]
}

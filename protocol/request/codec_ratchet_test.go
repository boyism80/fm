package request

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// Empty Serialize/Deserialize counts may only go down.
const (
	maxRequestSerializeEmpty    = 59
	maxRequestDeserializeEmpty  = 6
	maxResponseSerializeEmpty   = 1
	maxResponseDeserializeEmpty = 177
)

func TestEmptyCodecDoesNotGrow(t *testing.T) {
	reqSer, reqDe := countEmptyCodec(t, ".")
	respSer, respDe := countEmptyCodec(t, "../response")
	if reqSer > maxRequestSerializeEmpty || reqDe > maxRequestDeserializeEmpty {
		t.Fatalf("request empty serialize=%d deserialize=%d", reqSer, reqDe)
	}
	if respSer > maxResponseSerializeEmpty || respDe > maxResponseDeserializeEmpty {
		t.Fatalf("response empty serialize=%d deserialize=%d", respSer, respDe)
	}
}

func countEmptyCodec(t *testing.T, dir string) (serialize int, deserialize int) {
	t.Helper()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || fn.Body == nil {
					return true
				}
				if !emptyBody(fn.Body) {
					return true
				}
				switch fn.Name.Name {
				case "Serialize":
					serialize++
				case "Deserialize":
					deserialize++
				}
				return true
			})
		}
	}
	return serialize, deserialize
}

func emptyBody(body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return true
	}
	if len(body.List) != 1 {
		return false
	}
	ret, ok := body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	id, ok := ret.Results[0].(*ast.Ident)
	return ok && id.Name == "nil"
}

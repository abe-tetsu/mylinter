package myanalyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "myanalyzer is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "myanalyzer",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ForStmt:
			// pass.Reportf(n.Pos(), "for found")
			findLoopVar(pass, n)
		}
	})

	return nil, nil
}

func findLoopVar(pass *analysis.Pass, forstmt *ast.ForStmt) {
	assignStmt, ok := forstmt.Init.(*ast.AssignStmt)
	if !ok {
		return
	}

	// 左に定義されている変数がない場合は返す
	if len(assignStmt.Lhs) == 0 {
		return
	}

	// ループ変数を取得
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return
	}
	// ループ変数のスコープを取得
	obj := pass.TypesInfo.ObjectOf(ident)

	// For 文のスコープを取得
	forStmtScope, ok := pass.TypesInfo.Scopes[forstmt]
	if !ok {
		return
	}

	// ループ変数の一個上のスコープとFor文のスコープが一致しない場合は返す
	if forStmtScope != obj.Parent() {
		return
	}

	pass.Reportf(assignStmt.Pos(), "%v found", ident)
	findPointerOfLoopVar(pass, forStmtScope, forstmt.Body)
}

func findPointerOfLoopVar(pass *analysis.Pass, forstmtScope *types.Scope, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.UnaryExpr:
			// & じゃなかったら返す
			// TODO: これがなくてもテストが通るのがおかしい
			if n.Op != token.AND {
				return true
			}

			// x -> &の引数
			x, ok := n.X.(*ast.Ident)
			if !ok {
				return false
			}

			obj := pass.TypesInfo.ObjectOf(x)

			// xが宣言されている場所が、ループ変数と一致するときレポート
			if obj.Parent() == forstmtScope {
				pass.Reportf(n.Pos(), "unary expr found")
			}
		}
		return true
	})
}

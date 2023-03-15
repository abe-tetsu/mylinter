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

func findLoopVar(pass *analysis.Pass, forStmt *ast.ForStmt) {
	// For文の初期化文を取得
	assignStmt, ok := forStmt.Init.(*ast.AssignStmt)
	if !ok {
		return
	}

	// 左に定義されている変数がない場合は返す
	// foreach文の場合は対象外
	if len(assignStmt.Lhs) != 1 {
		return
	}

	// ループ変数を取得
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return
	}
	// ループ変数のスコープを取得
	obj := pass.TypesInfo.ObjectOf(ident)

	// For文のスコープを取得
	forStmtScope, ok := pass.TypesInfo.Scopes[forStmt]
	if !ok {
		return
	}

	// ループ変数の一個上のスコープとFor文のスコープが一致しない場合は返す
	if forStmtScope != obj.Parent() {
		return
	}

	pass.Reportf(assignStmt.Pos(), "%v found", ident)

	// for文のループ変数がポインタになっているかチェック
	findPointerOfLoopVar(pass, forStmtScope, forStmt.Body)
}

func findPointerOfLoopVar(pass *analysis.Pass, forStmtScope *types.Scope, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		// 関数呼び出しに対して処理
		case *ast.CallExpr:
			// 引数の数が0の場合は返す
			if len(n.Args) == 0 {
				return false
			}

			// 全ての引数に対して処理
			for _, arg := range n.Args {
				// 関数呼び出しを取得
				funcCall, ok := arg.(*ast.UnaryExpr)
				if !ok {
					continue
				}

				// & じゃなかったら返す
				// TODO: これがなくてもテストが通るのがおかしい
				if funcCall.Op != token.AND {
					continue
				}

				// x -> &の引数
				x, ok := funcCall.X.(*ast.Ident)
				if !ok {
					continue
				}

				obj := pass.TypesInfo.ObjectOf(x)

				// xが宣言されている場所が、ループ変数と一致するときレポート
				if obj.Parent() == forStmtScope {
					pass.Reportf(n.Pos(), "unary expr found")
				}
			}
		}
		return true
	})
}

package rewriter

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/tdakkota/gomacro"
)

type Runner struct {
	errorPrinter
	failed bool
}

func NewRunner(fset *token.FileSet) *Runner {
	return &Runner{
		errorPrinter: errorPrinter{NewErrorCallback(fset)},
	}
}

func (r *Runner) setReport(context *macro.Context) {
	context.Report = func(report macro.Report) {
		r.Reportf(report.Pos, report.Message)
	}
}

func (r *Runner) Run(handler macro.Handler, context macro.Context, node ast.Node) ast.Node {
	r.setReport(&context)

	return astutil.Apply(node, func(cursor *astutil.Cursor) bool {
		context.Pre = true
		context.Cursor = cursor

		if v := cursor.Node(); v != nil {
			err := handler.Handle(context, v)
			if err != nil {
				return false
			}
		}
		return true
	}, r.post(handler, context))
}

func (r *Runner) post(handler macro.Handler, context macro.Context) astutil.ApplyFunc {
	r.setReport(&context)

	return func(cursor *astutil.Cursor) bool {
		context.Cursor = cursor

		if v := cursor.Node(); v != nil {
			err := handler.Handle(context, v)
			if err != nil {
				r.err(v.Pos(), err)
				r.failed = true
				return false
			}
		}

		return true
	}
}

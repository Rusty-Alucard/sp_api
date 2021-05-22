// +build tools

// ソースでimportしないが、ツールとして利用するものを管理するためのファイル
// 通常のビルドには含めない

package tools

import (
	_ "github.com/99designs/gqlgen/cmd"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

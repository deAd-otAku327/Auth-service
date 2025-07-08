package repoerrors

import "errors"

// Repository errors.
var (
	ErrTransactionBegin  = errors.New("repo: tx begin failed")
	ErrTransactionCommit = errors.New("repo: tx commit failed")
	ErrQueryBuilding     = errors.New("repo: query building failed")
	ErrQueryExec         = errors.New("repo: query execution failed")
)

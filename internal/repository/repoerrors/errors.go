package repoerrors

import "errors"

// Repository errors.
var (
	ErrQueryBuildingFailed = errors.New("repo: query building failed")
	ErrQueryExecFailed     = errors.New("repo: query execution failed")
)

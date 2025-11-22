package kurohelpercore

import "errors"

var (
	// search no content for response
	ErrSearchNoContent = errors.New("search: no content for response")
	// The remote server returns a non-200 response status code
	ErrStatusCodeAbnormal = errors.New("response: server returned an error status code")
	// rate limit
	ErrRateLimit = errors.New("rate limit: rate limit, quota exhausted")
	// cache lost or expired
	ErrCacheLost = errors.New("cache: cache lost or expired")

	//ymgal invalid access token(401)
	ErrYmgalInvalidAccessToken = errors.New("ymgal: invalid access token or other 401 error")
	// trying to use bangumi character list search
	ErrBangumiCharacterListSearchNotSupported = errors.New("bangumi: character list search is not currently supported")
)

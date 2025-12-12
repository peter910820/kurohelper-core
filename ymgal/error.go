package ymgal

import (
	"errors"
	"fmt"
)

type (
	// api failed
	ErrAPIFailed struct {
		Code int
	}
)

var (
	//invalid access token(401)
	ErrInvalidAccessToken = errors.New("ymgal: invalid access token or other 401 error")
	//config invalid
	ErrCfgInvalid = errors.New("ymgal: invalid config, you need to call the init method to set up the default values")
)

func (e ErrAPIFailed) Error() string {
	return fmt.Sprintf("ymgal: api failed, error code=%d", e.Code)
}

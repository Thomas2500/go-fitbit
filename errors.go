package fitbit

import "errors"

var errExpiredToken = errors.New("expired token")
var errTokenChangeNotDefined = errors.New("tokenchange function is not defined")

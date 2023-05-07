package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/stuneak/simplebank/util"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {

	currecy, ok := fieldlevel.Field().Interface().(string)

	if ok {
		return util.IsSupportedCurrency(currecy)
	}

	return false
}

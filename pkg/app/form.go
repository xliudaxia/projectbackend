package app

import "github.com/gin-gonic/gin"

type ValidError struct {
	Message string
}

type ValidErrors []*ValidError

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		errs = append(errs, &ValidError{
			Message: err.Error(),
		})

		return false, errs
	}

	return true, nil
}

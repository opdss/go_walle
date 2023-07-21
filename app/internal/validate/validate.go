package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-walle/app/model/field"
)

var validateFuncs map[string]func(fl validator.FieldLevel) bool

func init() {
	validateFuncs = make(map[string]func(fl validator.FieldLevel) bool)
	validateFuncs["status"] = func(fl validator.FieldLevel) bool {
		t, ok := fl.Field().Interface().(field.Status)
		if ok {
			t1 := int(t)
			if t1 == field.StatusEnable || t1 == field.StatusDisable {
				return true
			}
		}
		return false
	}

	//validateFuncs["userid"] = func(fl validator.FieldLevel) bool {
	//	params := strings.Kv2MapString(fl.Param())
	//	v, ok := fl.Field().Interface().(int64)
	//	if ok {
	//		if v == 0 {
	//			return false
	//		}
	//		m, err := user.NewService().Detail(v)
	//		if err != nil {
	//			return false
	//		}
	//		if st, ok := params["status"]; ok {
	//			if st == "" {
	//				return false
	//			}
	//			status, err := strconv.Atoi(st)
	//			if err != nil {
	//				return false
	//			}
	//			if int(m.Status) == status {
	//				return true
	//			} else {
	//				return false
	//			}
	//		}
	//		return true
	//	}
	//	return false
	//}
}

func RegisterValidation() error {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for k, fn := range validateFuncs {
			err := validate.RegisterValidation(k, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

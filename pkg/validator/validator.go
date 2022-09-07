package validator

import (
	"github.com/abdivasiyev/project_template/pkg/translator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
	"sync"
)

type TranslatableValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &TranslatableValidator{}

func (v *TranslatableValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *TranslatableValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *TranslatableValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		_ = enTranslations.RegisterDefaultTranslations(v.validate, translator.Get())
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

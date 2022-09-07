package translator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var _instance ut.Translator = nil

func Get() ut.Translator {
	if _instance == nil {
		translator := en.New()
		uni := ut.New(translator, translator)

		_instance, _ = uni.GetTranslator("en")
	}

	return _instance
}

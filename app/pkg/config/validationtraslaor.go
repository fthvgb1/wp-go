package config

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

var enT ut.Translator
var zhT ut.Translator

func GetZh() ut.Translator {
	return zhT
}
func GetEn() ut.Translator {
	return enT
}

func InitTrans() error {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		ens := en.New()
		uni := ut.New(ens, zh.New(), ens)
		zhT, _ = uni.GetTranslator("zh")
		enT, _ = uni.GetTranslator("en")
		err := enTrans.RegisterDefaultTranslations(validate, enT)
		if err != nil {
			return err
		}
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("label")
		})
		err = zhTrans.RegisterDefaultTranslations(validate, zhT)
		if err != nil {
			return err
		}
	}
	return nil
}

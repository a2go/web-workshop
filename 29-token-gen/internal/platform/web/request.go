package web

import (
	"encoding/json"
	"errors"
	"net/http"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// v holds the settings and caches for validating request struct values.
var v = validator.New()

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

func init() {

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, lang)
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return ErrorWithStatus(err, http.StatusBadRequest)
	}

	if err := v.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verr, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields []fieldError
		for field, msg := range verr.Translate(lang) {
			fields = append(
				fields,
				fieldError{Field: field, Error: msg},
			)
		}

		return &statusError{
			err:    errors.New("field validation error"),
			status: http.StatusBadRequest,
			fields: fields,
		}
	}

	return nil
}

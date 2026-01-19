package main

import (
	"io"
	"bytes"
	"net/http"
	"time"
	"fmt"
	"strings"

	"encoding/json"

	"errors"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)

}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	//func NewBuffer(buf []byte) *Buffer
	buf := new(bytes.Buffer)

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)

	if err != nil {
		err = fmt.Errorf("cannot write to buffer")
		app.serverError(w, r, err)
		return
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {

	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		//CSRFToken      : nosurf.Token(r),
	}

}

func (app *application) decodePostForm(r *http.Request, dst any) error {

	// ParseForm parses the raw query from the URL and updates r.Form
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}


func DecodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
    maxBytes := 1_048_576
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    err := decoder.Decode(destination)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError
        var invalidUnmarshalError *json.InvalidUnmarshalError
        var maxBytesError *http.MaxBytesError

        switch {
        case errors.As(err, &syntaxError):
            return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("body contains badly-formed JSON")

        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
            }

            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

        case errors.Is(err, io.EOF):
            return errors.New("body must not be empty")

        // Setting DisallowUnknownFields may trigger this error when an unknown field is found.
        // Since there is no specific error type for this case a string match is used.
        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            return fmt.Errorf("body contains unknown key %s", fieldName)

        case errors.As(err, &maxBytesError):
            return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

        case errors.As(err, &invalidUnmarshalError):
            panic(err)

        default:
            return fmt.Errorf("failed to parse JSON: %w", err)
        }
    }

    err = decoder.Decode(&struct{}{})
    if !errors.Is(err, io.EOF) {
        return errors.New("body must only contain a single JSON value")
    }

    return nil
}



func (app *application) isAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(isAuthenticatedContextKey).(bool)

	app.logger.Info("Is AUtheticated", "KEY", isAuth)
	if !ok {
		return false
	}
	return isAuth
}

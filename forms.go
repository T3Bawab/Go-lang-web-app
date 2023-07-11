package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) HasRequired(IdTags ...string) {
	for _, Idtag := range IdTags {
		value := f.Get(Idtag)

		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(Idtag, "This field can't be blank")
		}

	}

}

func (f *Form) HasValue(IdTag string, r *http.Request) bool {
	x := r.Form.Get(IdTag)

	return x != ""
}

func (f *Form) MinLength(IdTag string, lenght int, r *http.Request) bool {
	x := r.Form.Get(IdTag)

	if len(x) < lenght {
		f.Errors.AddError(IdTag, fmt.Sprintf("must be %d charcaters long or more ", lenght))
		return false

	}

	return true

}

func (f *Form) IsEmail(IdTag string) {
	if !govalidator.IsEmail(f.Get(IdTag)) {
		f.Errors.AddError(IdTag, "Invalid Email")
	}

}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

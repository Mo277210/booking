package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct, embeds a url.Values object (from net/url package) and adds an Errors field to hold form validation errors.
type Form struct {
url.Values
Errors errors

}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool{
	return len(f.Errors)==0
}

// New initializes a custom form struct
func New(data url.Values) *Form{
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}
//Has checks if a field is present in the form data
func (f *Form) Has(field string,r *http.Request) bool{
	x:=r.Form.Get(field)
	if x==""{
		f.Errors.Add(field,"This field cannot be blank")	
		return false
	}
	return true

}
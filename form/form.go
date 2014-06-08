package form

import (
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/gotcha/http"
	"reflect"
	"strings"
)

type FormHelper struct {
	Session   *http.Session
	Model     interface{}
	Values    map[string]string
	HasErrors bool
	Rules     map[string][]*Rule
	Errors    map[string]map[string]error
}

func New(session *http.Session, model interface{}) *FormHelper {
	fh := &FormHelper{
		Session:   session,
		Model:     model,
		Values:    make(map[string]string),
		HasErrors: false,
		Rules:     make(map[string][]*Rule),
		Errors:    make(map[string]map[string]error),
	}
	fh.parseRules()
	return fh
}

func (fh *FormHelper) parseRules() {
	s := reflect.TypeOf(fh.Model).Elem()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i).Name

		tag := s.Field(i).Tag
		if len(tag) > 0 {
			log.Trace("Model field [%s] has validation tag: %s", f, tag)

			s := strings.Split(string(tag), ";")
			for _, rule := range s {
				p := strings.SplitN(rule, ":", 2)
				if len(p) == 1 {
					rule = strings.TrimSpace(p[0])
					log.Trace("Found rule [%s]", rule)
				} else {
					rule = strings.TrimSpace(p[0])
					params := strings.TrimSpace(p[1])
					log.Trace("Found rule [%s] with parameters [%s]", rule, params)
					if _, ok := fh.Rules[f]; !ok {
						fh.Rules[f] = make([]*Rule, 0)
					}
					switch rule {
					case "minlength":
						fh.Rules[f] = append(fh.Rules[f], MinLength(params))
					case "maxlength":
						fh.Rules[f] = append(fh.Rules[f], MaxLength(params))
					}
				}
			}
		}
	}
}

func (fh *FormHelper) Populate(multipart bool) *FormHelper {
	// TODO nested form values

	t := reflect.TypeOf(fh.Model)
	v := reflect.ValueOf(fh.Model)
	s := t.Elem()
	w := v.Elem()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i).Name
		fl := strings.ToLower(f)
		log.Trace("Model field [%s] mapped to form field [%s]", f, fl)

		var val []string
		if multipart {
			val = fh.Session.Request.MultipartForm().Value[fl]
		} else {
			val = fh.Session.Request.PostForm()[fl]
		}
		if len(val) > 0 {
			fh.Values[f] = val[0]
			w.Field(i).SetString(val[0])
			log.Trace("Field [%s] had value [%s]", fl, fh.Values[fl])
		}
	}

	return fh
}

func (fh *FormHelper) Validate() *FormHelper {
	for k, f := range fh.Rules {
		log.Trace("Validating field [%s] with value [%s]", k, fh.Values[k])

		valid := true
		for _, r := range f {
			log.Trace("Executing rule [%s] with parameters [%s]", r.Name, r.Parameters)

			err := r.Function(fh.Values[k])

			if err != nil {
				log.Trace("Rule [%s] failed: %s", r.Name, err)
				fh.HasErrors = true
				valid = false
				if _, ok := fh.Errors[k]; !ok {
					fh.Errors[k] = make(map[string]error)
				}
				fh.Errors[k][r.Name] = err
			} else {
				log.Trace("Rule [%s] passed", r.Name)
			}
		}

		log.Trace("Field [%s] is valid: %t", k, valid)
	}

	return fh
}

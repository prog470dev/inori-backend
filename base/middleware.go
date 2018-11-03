package base

import (
	"net/http"
)

func middle(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// do something before f
		f(w, r)
		// do something after f
	}
}

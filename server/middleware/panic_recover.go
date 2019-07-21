package middleware

import (
	"log"
	"net/http"
	"runtime"
)

func PanicRecover(handler http.HandlerFunc) http.HandlerFunc {
	var buffer [1 << 16]byte
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stackLen := runtime.Stack(buffer[:], false)
				log.Printf("panic occur: %v, current stack trace(may not cover all): %s", r, buffer[:stackLen])
			}
		}()
		handler(w, r)
	}
}

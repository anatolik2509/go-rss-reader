package security

import "net/http"

type HttpHandler struct{}

func (*HttpHandler) SignIn(r *http.Request, w http.ResponseWriter) {

}

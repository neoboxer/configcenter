package internal

import (
	"net/http"
)

type internalServer struct {
}

func (s *internalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

package ayame

import (
	"fmt"
	"net/http"
)

func (s *Server) okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

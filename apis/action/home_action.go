package action

import (
	"net/http"

	"github.com/naufalziyad/RESTful-APIs/apis/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Started to this RESTful APIs")
}

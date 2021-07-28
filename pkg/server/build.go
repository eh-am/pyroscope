package server

import (
	"net/http"

	"github.com/pyroscope-io/pyroscope/pkg/build"
)

func (ctrl *Controller) buildHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(build.PrettyJSON()))

	if err != nil {
		ctrl.log.WithError(err).Error("Build Handler failed")
	}
}

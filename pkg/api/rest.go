package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/qubole/wireguard/pkg/wgclient"
	"bitbucket.org/qubole/wireguard/pkg/wgserver"
)

// REST apis
type REST struct {
	WGS *wgserver.Svc
	WGC *wgclient.Svc
}

// StatusHandler is for any http status code.
type StatusHandler struct {
	Err  error
	Code int
}

// ServeHTTP is http.Handler insterface implementation.
func (s StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeError(s.Err, s.Code, w)
}

// ClientGererateConfig returns wgclient config:
// Input:
// // {
// // 	"id": "5",
// // 	"public_key": "dhfjdbfjdbffg"
// // }
// Output:
// //   {
// //    "id": "5",
// //    "private_ip": "10.0.0.3",
// //    "ssh_authorized_keys": [
// //        "test1",
// //        "test2"
// //    ],
// //    "public_key": "dhfjdbfjdbffg"
// }
func (h *REST) ClientGererateConfig() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in wgclient.GenerateConfigInput

		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			writeError(fmt.Errorf("wgclient:create:%v", err), http.StatusBadRequest, w)
			return
		}

		out, err := h.WGC.GenerateConfig(r.Context(), &in)
		if err != nil {
			writeError(fmt.Errorf("wgclient:create:%v", err), http.StatusBadRequest, w)
			return
		}

		writeRespone(out, w)
	})
}

// writeError writes error on ResponseWriter
func writeError(err error, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// writeError writes error on ResponseWriter
func writeRespone(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

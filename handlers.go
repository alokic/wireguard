package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/qubole/wireguard/pkg/wgclient"
	"bitbucket.org/qubole/wireguard/pkg/wgserver"
)

type handlers struct {
	cfg *Config
	s   *wgserver.Svc
	c   *wgclient.Svc
}

// getClientProfile returns wgclient specifc Profile:
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
func (h *handlers) getClientProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in wgclient.CreateInput

		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			writeError(fmt.Errorf("wgclient:create:%v", err), http.StatusBadRequest, w)
			return
		}

		out, err := h.c.Create(r.Context(), &in)
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

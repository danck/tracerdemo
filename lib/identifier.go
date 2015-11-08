package tracerdemo

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
)

// Identifier is a decorator for HandlerFuncs that attaches a GUID to each
// incoming request
func Identifier(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// generate a uuid based on timestamp and hwaddress
		//	note: if this runs in a docker container
		//	With the current implementation hwaddresses are not necessarily
		//	unique if the container's ip address isn't:
		//	https://github.com/docker/libnetwork/blob/master/netutils/utils.go#L115-L132 (08.11.2015)
		uuid := uuid.NewV1()
		r.Header.Set("consolidation-id", fmt.Sprint(uuid))
		fn(w, r)
	})
}

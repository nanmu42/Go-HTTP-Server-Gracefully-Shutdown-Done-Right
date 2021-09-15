package downright

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func SlowHandler(sleepSeconds int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = io.WriteString(w, "only accept GET request")
		}

		id := randString(4)

		log.Printf("[id: %s] new request from %s, sleep for %d seconds before response...\n", id, req.RemoteAddr, sleepSeconds)

		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		log.Printf("[id: %s] sleep done, now responding...\n", id)

		var err error
		_, err = fmt.Fprintf(w, "Hi, your request id is %s", id)
		if err != nil {
			log.Printf("[id: %s] error while wring response: %s\n", id, err)
		} else {
			log.Printf("[id: %s] responded successfully.\n", id)
		}
	})
}

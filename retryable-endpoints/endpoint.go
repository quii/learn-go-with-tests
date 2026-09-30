package retryableendpoints

import (
	"encoding/json"
	"errors"
	"net/http"
)

func RetryableEndpoint(topUps TopUps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request TopUpRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := topUps.Apply(r.Context(), r.Header.Get("Idempotency-Key"), request)
		if errors.Is(err, ErrMissingKey) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "could not top up account", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Once writing starts, an encoding/write error cannot change the status.
		// The completed result remains available for a retry.
		_ = json.NewEncoder(w).Encode(result)
	})
}

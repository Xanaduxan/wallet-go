package handlers

import "net/http"

func Me(w http.ResponseWriter, r *http.Request) {

	email := r.Context().Value("email").(string)

	writeJSON(w, http.StatusOK, map[string]string{
		"email": email,
	})
}

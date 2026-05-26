package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev mode."))
		return
	}

	if err := cfg.dbQueries.DeleteAllUsers(req.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to perform the deletion operation", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

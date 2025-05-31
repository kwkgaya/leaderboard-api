package api

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/leaderboard", leaderboardHandler)
	return mux
}


func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	leaderboard := []string{"Alice", "Bob", "Charlie"}
	jsonData, err := json.Marshal(leaderboard)
    if err != nil {
        fmt.Println("Error marshalling array:", err)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"leaderboard": %v}`, string(jsonData))
}
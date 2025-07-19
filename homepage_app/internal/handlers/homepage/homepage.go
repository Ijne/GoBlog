package homepage

import "net/http"

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Homepage!"))
}

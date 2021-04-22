package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	lb := &LeaderBoard{}

	http.HandleFunc("/usercount", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			count := lb.UserCount()
			fmt.Fprintf(w, `{"count":%v}`, count)
		default:
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		userId := strings.TrimPrefix(r.URL.EscapedPath(), "/users/")
		switch r.Method {
		case http.MethodGet:
			user := lb.GetUser(userId)
			data, err := json.Marshal(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPut:
			score, err := getQueryParamInt(r, "score")
			if err != nil {
				http.Error(w, "cannot parse score", http.StatusBadRequest)
				return
			}

			lb.SetUser(userId, score)
			user := lb.GetUser(userId)
			data, err := json.Marshal(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/ranks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rank, err := getQueryParamInt(r, "rank")
			if err != nil {
				http.Error(w, "cannot parse rank", http.StatusBadRequest)
				return
			}

			count, err := getQueryParamInt(r, "count")
			if err != nil {
				http.Error(w, "cannot parse count", http.StatusBadRequest)
				return
			}

			users := lb.GetRanks(rank, count)
			data, err := json.Marshal(users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.NotFound(w, r)
		}
	})

	fmt.Println("listen..")
	http.ListenAndServe(":8080", nil)
}

func getQueryParamInt(r *http.Request, name string) (int, error) {
	s := strings.TrimSpace(r.URL.Query().Get(name))
	if s == "" {
		return 0, errors.New("empty")
	}

	return strconv.Atoi(s)
}

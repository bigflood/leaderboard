package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Router struct {
	lb *LeaderBoard
}

func NewRouter(lb *LeaderBoard) *Router {
	return &Router{lb: lb}
}

func (router *Router) Setup() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/usercount", router.HandleUserCount)
	mux.HandleFunc("/users/", router.HandleUsers)
	mux.HandleFunc("/ranks", router.HandleRanks)
	return mux
}

func (router *Router) HandleUserCount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		count := router.lb.UserCount()
		fmt.Fprintf(w, `{"count":%v}`, count)
	default:
		http.NotFound(w, r)
	}
}

func (router *Router) HandleUsers(w http.ResponseWriter, r *http.Request) {
	userId := strings.TrimPrefix(r.URL.EscapedPath(), "/users/")
	switch r.Method {
	case http.MethodGet:
		user := router.lb.GetUser(userId)
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

		router.lb.SetUser(userId, score)
		user := router.lb.GetUser(userId)
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
}

func (router *Router) HandleRanks(w http.ResponseWriter, r *http.Request) {
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

		users := router.lb.GetRanks(rank, count)
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
}

func getQueryParamInt(r *http.Request, name string) (int, error) {
	s := strings.TrimSpace(r.URL.Query().Get(name))
	if s == "" {
		return 0, errors.New("empty")
	}

	return strconv.Atoi(s)
}

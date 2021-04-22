package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bigflood/leaderboard/leaderboard"
)

type HttpHandler struct {
	lb *leaderboard.LeaderBoard
}

func New(lb *leaderboard.LeaderBoard) *HttpHandler {
	return &HttpHandler{lb: lb}
}

func (handler *HttpHandler) Setup() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/usercount", handler.HandleUserCount)
	mux.HandleFunc("/users/", handler.HandleUsers)
	mux.HandleFunc("/ranks", handler.HandleRanks)
	return mux
}

func (handler *HttpHandler) HandleUserCount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		count := handler.lb.UserCount()
		fmt.Fprintf(w, `{"count":%v}`, count)
	default:
		http.NotFound(w, r)
	}
}

func (handler *HttpHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	userId := strings.TrimPrefix(r.URL.EscapedPath(), "/users/")
	switch r.Method {
	case http.MethodGet:
		user := handler.lb.GetUser(userId)
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

		handler.lb.SetUser(userId, score)
		user := handler.lb.GetUser(userId)
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

func (handler *HttpHandler) HandleRanks(w http.ResponseWriter, r *http.Request) {
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

		users := handler.lb.GetRanks(rank, count)
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

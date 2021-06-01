package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
)

// Handler for http requests
type Handler struct {
	Router *mux.Router
	cfg    config.Config
	Events chan node.Node
	t      *temper.Temper
	l      *log.Entry
}

// New http handler
func New(cfg config.Config, l *log.Entry, t *temper.Temper) *Handler {
	e := make(chan node.Node)
	h := Handler{mux.NewRouter(), cfg, e, t, l}
	return &h
}

// RegisterEventRoute for a node event endpoint
func (h *Handler) RegisterEventRoute() {
	h.Router.HandleFunc("/events/", h.eventHandler)
}

// RegisterAPIRoutes for a node event endpoint
func (h *Handler) RegisterAPIRoutes() {
	h.Router.HandleFunc("/api/nodes/{node}/tasks/{task}", h.temperHandler).Methods("POST")
	h.Router.HandleFunc("/api/nodes", h.nodeListHandler).Methods("GET")
	if h.t != nil {
		h.Router.HandleFunc("/api/nodes/webhook", h.webhookHandler).Methods("POST")
	}
}

func (h *Handler) nodeListHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(h.t.GetNodes()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) temperHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := vars["node"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.execTasks(n, r.URL, r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "node: %v\n", n)
}

func (h *Handler) webhookHandler(w http.ResponseWriter, r *http.Request) {
	wb := webhookBody{}
	if err := json.NewDecoder(r.Body).Decode(&wb); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	n, _ := node.New("test", h.cfg)
	h.t.AddNodes([]*node.Node{n})
}

func (h *Handler) eventHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	defer r.Body.Close()
	by, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bs := string(by)
	fmt.Println(p[2], bs)
}

func (h *Handler) execTasks(n string, u *url.URL, ctx context.Context) (err error) {
	node, err := node.New(n, h.cfg)
	if err != nil {
		return
	}
	vals := u.Query()["task"]
	for _, v := range vals {
		switch v {
		case "sync_netbox":
			node.Update()
		case "cablecheck":
		}
	}
	return
}

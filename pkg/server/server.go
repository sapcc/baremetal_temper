package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
)

// Handler for http requests
type Handler struct {
	router *mux.Router
	cfg    config.Config
	Events chan model.Node
	l      *log.Entry
}

// New http handler
func New(cfg config.Config, l *log.Entry) *Handler {
	e := make(chan model.Node)
	h := Handler{mux.NewRouter(), cfg, e, l}
	return &h
}

// RegisterEventRoutes for a node event endpoint
func (h *Handler) RegisterEventRoute() {
	h.router.HandleFunc("baremetal_temper/events/", h.eventHandler)
}

// RegisterEventRoutes for a node event endpoint
func (h *Handler) RegisterTemperRoutes() {
	h.router.HandleFunc("baremetal_temper/temper/{node}", h.temperHandler)
}

func (h *Handler) temperHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := vars["node"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "node: %v\n", n)
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

func (h *Handler) execTasks(n string, u *url.URL) (v string, err error) {
	node := model.Node{Name: n}
	t := temper.New(h.cfg)
	c, err := t.GetClients(n)
	vals := u.Query()["task"]
	for _, v := range vals {
		switch v {
		case "sync_netbox":
			c.Netbox.LoadInterfaces(&node)
		case "cablecheck":
		}
	}
	return
}

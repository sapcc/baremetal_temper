package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sapcc/baremetal_temper/pkg/model"
)

// Handler for http requests
type Handler struct {
	mux    *http.ServeMux
	Events chan model.Node
}

// New http handler
func New() *Handler {
	e := make(chan model.Node)
	mux := http.NewServeMux()
	h := Handler{mux, e}

	return &h
}

// RegisterEventRoutes for a node event endpoint
func (h *Handler) RegisterEventRoute(n *model.Node) {
	path := fmt.Sprintf("baremetal_temper/events/%s", n.Name)
	h.mux.HandleFunc(path, h.eventHandler)
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

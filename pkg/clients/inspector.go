package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type NodeAlreadyExists struct {
	Err string
}

func (n *NodeAlreadyExists) Error() string {
	return n.Err
}

type InspectorErr struct {
	Error ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type InspectorClient struct {
	log  *log.Entry
	host string
}

func NewInspectorClient(cfg config.Config, ctxLogger *log.Entry) *InspectorClient {
	return &InspectorClient{
		log:  ctxLogger,
		host: cfg.Inspector.Host,
	}
}

func (i InspectorClient) CreateIronicNode(in *model.IronicNode) (err error) {
	i.log.Info("calling inspector api for node creation")
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("http://%s", i.host))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "/v1/continue")
	db, err := json.Marshal(in.InspectionData)
	if err != nil {
		return
	}
	log.Debugf("calling (%s) with data: %s", u.String(), string(db))
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(db))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		ierr := &InspectorErr{}
		if err = json.Unmarshal(bodyBytes, ierr); err != nil {
			return fmt.Errorf("could not create node")
		}
		if strings.Contains(ierr.Error.Message, "already exists, uuid") {
			return &NodeAlreadyExists{}
		}
		return fmt.Errorf(ierr.Error.Message)
	}

	if err = json.Unmarshal(bodyBytes, in); err != nil {
		return
	}
	name := strings.Split(in.InspectionData.Inventory.BmcAddress, ".")
	node := strings.Split(name[0], "-")
	nodeName := strings.Replace(node[0], "r", "", 1)
	in.Name = fmt.Sprintf("%s-%s", nodeName, node[1])
	return
}

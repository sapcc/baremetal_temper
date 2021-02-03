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

//NodeAlreadyExists custom error
type NodeAlreadyExists struct {
	Err string
}

func (n *NodeAlreadyExists) Error() string {
	return n.Err
}

//InspectorErr custom error struct for inspector callback errors
type InspectorErr struct {
	Error ErrorMessage `json:"error"`
}

//ErrorMessage message struct for InspectorErr
type ErrorMessage struct {
	Message string `json:"message"`
}

//InspectorClient is
type InspectorClient struct {
	log  *log.Entry
	host string
}

//NewInspectorClient creates a ironic-inspector client
func NewInspectorClient(cfg config.Config, ctxLogger *log.Entry) *InspectorClient {
	return &InspectorClient{
		log:  ctxLogger,
		host: cfg.Inspector.Host,
	}
}

//Create creates a new ironic node based on the provided ironic model
func (i InspectorClient) Create(in *model.IronicNode) (err error) {
	i.log.Debug("calling inspector api for node creation")
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

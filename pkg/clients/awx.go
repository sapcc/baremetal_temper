package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

type AwxClient struct {
	cfg config.Config
	log *log.Entry
}

type launchBody struct {
	Inventory int    `json:"inventory"`
	Limit     string `json:"limit"`
}

func (a AwxClient) ExecTemplates(host string) (err error) {
	lb := launchBody{
		Inventory: 94,
		Limit:     host,
	}
	err, j := a.execTemplate(lb, "")
	if err != nil {
		return
	}
	if err = a.checkJobStatus(j); err != nil {
		return
	}
	return
}

func (a AwxClient) execTemplate(l launchBody, temp string) (err error, job string) {
	cfg := a.cfg.Awx
	u := fmt.Sprintf("%s/job_templates/%s/launch", cfg.Host, temp)
	b, err := json.Marshal(&l)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(b))
	req.SetBasicAuth(cfg.User, cfg.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error exec template"), job
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var r map[string]json.RawMessage
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		return
	}
	err = json.Unmarshal(r["job"], &job)
	return
}

func (a AwxClient) checkJobStatus(job string) (err error) {
	cfg := a.cfg.Awx
	u := fmt.Sprintf("%s/jobs/%s/stdout/", cfg.Host, job)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(cfg.User, cfg.Password)
	cf := wait.ConditionFunc(func() (bool, error) {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return true, err
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		var r map[string]json.RawMessage
		err = json.Unmarshal(bodyBytes, &r)
		var status string
		var failed bool
		err = json.Unmarshal(r["status"], &status)
		err = json.Unmarshal(r["failed"], &failed)
		if status != "successful" {
			return false, nil
		}
		if failed {
			return true, fmt.Errorf("job exec failed")
		}
		return true, nil
	})
	if err = wait.Poll(5*time.Second, 20*time.Minute, cf); err != nil {
		return
	}

	return
}

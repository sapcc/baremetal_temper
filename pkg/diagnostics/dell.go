package diagnostics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Dell struct {
	client *gofish.APIClient
	log    *log.Entry
}

type iDracDiagnostics struct {
	RebootJobType string `json:"RebootJobType,omitempty"`
	RunMode       string `json:"RunMode,omitempty"`
	ShareType     string `json:"ShareType,omitempty"`
}

type iDRACJob struct {
	CompletionTime string `json:"CompletionTime"`
	JobType        string `json:"JobType"`
	JobState       string `json:"JobState"`
	ID             string `json:"Id"`
}

func NewDell(c *gofish.APIClient, l *log.Entry) Diagnostics {
	return &Dell{c, l}
}

func (d Dell) Run() (err error) {
	payload := iDracDiagnostics{RebootJobType: "GracefulRebootWithForcedShutdown", RunMode: "Express"} //Extended
	resp, err := d.client.Post("/redfish/v1/Dell/Managers/iDRAC.Embedded.1/DellLCService/Actions/DellLCService.RunePSADiagnostics", payload)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("run remote diags not successful")
	}
	loc := resp.Header.Get("Location")
	locs := strings.Split(loc, "/")
	jobID := locs[len(locs)-1]

	cf := wait.ConditionFunc(func() (bool, error) {
		j, err := d.getJobByID(jobID)
		if err != nil {
			log.Errorf("Error loading diagnostics job info: %s", err.Error())
			return false, err
		}
		if j.JobState == "Completed" {
			log.Debug("diagnostics job completed")
			return true, nil
		}
		log.Debugf("waiting for diagnostics job to be completed. state: %s", j.JobState)
		return false, nil
	})

	if err = wait.Poll(60*time.Second, 120*time.Minute, cf); err != nil {
		return err
	}

	res, err := d.getDiagnosticsResult()
	passed := true
	for r, i := range res {
		if i < 1 {
			log.Errorf("diagnostic test did not pass: %s", r)
			passed = false
		}
	}
	if !passed {
		return fmt.Errorf("diagnostic tests did not pass")
	}
	return
}

func (d Dell) getJobByID(id string) (j iDRACJob, err error) {
	resp, err := d.client.Get("/redfish/v1/Managers/iDRAC.Embedded.1/Jobs/" + id)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return
	}
	return
}

func (d Dell) getDiagnosticsResult() (results map[string]int, err error) {
	var rgx = regexp.MustCompile(`\*\*(.*?)\*\*`)
	var test string
	results = make(map[string]int)

	payload := iDracDiagnostics{ShareType: "Local"}
	resp, err := d.client.Post("/redfish/v1/Dell/Managers/iDRAC.Embedded.1/DellLCService/Actions/DellLCService.ExportePSADiagnosticsResult", payload)
	if err != nil {
		return
	}

	resp, err = d.client.Get(resp.Header.Get("Location"))
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "**") {
			rs := rgx.FindStringSubmatch(scanner.Text())
			test = rs[1]
			results[rs[1]] = 0
		}

		if strings.Contains(scanner.Text(), "Test Results :") {
			r := strings.Split(scanner.Text(), " : ")

			if _, ok := results[test]; ok {
				switch r[1] {
				case "Pass":
					results[test] = 2
				case "Warning":
					results[test] = 1
				default:
					results[test] = 0
				}
			}
		}
	}

	return
}

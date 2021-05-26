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

type DellClient struct {
	client *gofish.APIClient
	gCfg   gofish.ClientConfig
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

type iDRACError struct {
	Error iDRACErrorInfo `json:"error"`
}

type iDRACErrorInfo struct {
	Message []iDRACErrorMsg `json:"@Message.ExtendedInfo"`
}

type iDRACErrorMsg struct {
	Message string `json:"Message"`
}

type iDRACJobList struct {
	Members []iDRACJobMember `json:"Members"`
}

type iDRACJobMember struct {
	URL string `json:"@odata.id"`
}

func NewDellClient(gCfg gofish.ClientConfig, log *log.Entry) (c *DellClient) {
	return &DellClient{
		gCfg: gCfg,
		log:  log,
	}
}

func (d DellClient) Run() (err error) {
	client, err := gofish.Connect(d.gCfg)
	defer client.Logout()
	if err != nil {
		return
	}
	d.client = client
	payload := iDracDiagnostics{RebootJobType: "GracefulRebootWithForcedShutdown", RunMode: "Extended"} //Express
	resp, err := client.Post("/redfish/v1/Dell/Managers/iDRAC.Embedded.1/DellLCService/Actions/DellLCService.RunePSADiagnostics", payload)
	var jobID string
	if err != nil {
		idracErr := &iDRACError{}
		if err = json.Unmarshal([]byte(strings.Replace(err.Error(), "400: ", "", 1)), idracErr); err != nil {
			return
		}
		if idracErr.Error.Message[0].Message == "A Remote Diagnostic (ePSA) job already exists." {
			d.log.Errorf("remote diags already running")
			job, err := d.findDiagnosticJob()
			if err != nil {
				return err
			}
			jobID = job.ID
		} else {
			return fmt.Errorf(idracErr.Error.Message[0].Message)
		}

	} else if resp.StatusCode != 202 {
		return fmt.Errorf("run remote diags not successful")
	} else {
		loc := resp.Header.Get("Location")
		locs := strings.Split(loc, "/")
		jobID = locs[len(locs)-1]
	}
	jobStartFailed := false
	cf := wait.ConditionFunc(func() (bool, error) {
		j, err := d.getJobByID(jobID)
		if err != nil {
			d.log.Errorf("Error loading diagnostics job info: %s", err.Error())
			return false, err
		}
		if j.JobState == "Completed" {
			d.log.Debug("diagnostics job completed")
			return true, nil
		}
		if j.JobState == "" {
			d.log.Debug("waiting for diagnostics job to be completed. state: unavailable")
			jobStartFailed = true
			return true, nil
		}
		d.log.Debugf("waiting for diagnostics job to be completed. state: %s", j.JobState)
		return false, nil
	})

	if err = wait.Poll(60*time.Second, 240*time.Minute, cf); err != nil {
		return err
	}

	if jobStartFailed {
		return d.Run()
	}

	res, err := d.getDiagnosticsResult()
	passed := true
	for r, i := range res {
		if i < 1 {
			d.log.Errorf("diagnostic test did not pass: %s", r)
			passed = false
		}
	}
	if !passed {
		return fmt.Errorf("diagnostic tests did not pass")
	}
	return
}

func (d DellClient) getJobByID(id string) (j iDRACJob, err error) {
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

func (d DellClient) findDiagnosticJob() (j iDRACJob, err error) {
	resp, err := d.client.Get("/redfish/v1/Managers/iDRAC.Embedded.1/Jobs")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	jl := iDRACJobList{}
	err = json.NewDecoder(resp.Body).Decode(&jl)
	if err != nil {
		return
	}
	for _, m := range jl.Members {
		s := strings.Split(m.URL, "/")
		j, err := d.getJobByID(s[len(s)-1])
		if err != nil {
			return j, err
		}
		if j.JobType == "RemoteDiagnostics" && j.JobState == "Running" {
			d.log.Debugf("found diagnostics job %s", j.ID)
			return j, nil
		}
	}
	return
}

func (d DellClient) getDiagnosticsResult() (results map[string]int, err error) {
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

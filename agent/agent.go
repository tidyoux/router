package agent

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/tidyoux/router/agent/config"
	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/handler"
	log "github.com/sirupsen/logrus"
	"github.com/tidyoux/goutils"
)

type Agent struct {
	cfg    *config.Config
	client *Client
}

func New(cfg *config.Config) *Agent {
	return &Agent{
		cfg:    cfg,
		client: NewClient(cfg.URL),
	}
}

func (a *Agent) Name() string { return "agent" }

func (a *Agent) Init() error {
	return a.client.Login(a.cfg.WorkerID, a.cfg.WorkerKey)
}

func (a *Agent) Work() {
	resp, err := a.client.ListTask()
	if err != nil {
		if err.Error() == handler.ErrInvalidAgentToken.Error() {
			err = a.Init()
			if err == nil {
				return
			}
		}

		log.Errorf("list task failed, %v", err)
		return
	}

	for _, task := range resp.Tasks {
		success, detail, err := a.process(task)
		if err != nil {
			log.Errorf("process task (id: %d, params: %s) failed, %v",
				task.ID, task.Params, err)
			return
		}

		err = a.client.FinishTask(task.ID, success, detail)
		if err != nil {
			log.Errorf("finish task (id: %d, params: %s) failed, %v",
				task.ID, task.Params, err)
		} else {
			log.Infof("finish process task (id: %d, params: %s)",
				task.ID, task.Params)
		}
	}
}

func (a *Agent) process(task *types.Task) (bool, string, error) {
	log.Infof("start process task (id: %d, params: %s, process: %d, status: %d, detail: %s)",
		task.ID, task.Params, task.Progress, task.Status, task.Detail)

	if task.Status == types.TaskRecord {
		err := a.client.AcceptTask(task.ID)
		if err != nil {
			return false, "", fmt.Errorf("accept task failed, %v", err)
		}
	}

	if len(task.Params) == 0 {
		return false, "invalid task params, params can't be empty", nil
	}

	params := strings.Split(task.Params, " ")
	taskName := params[0]
	taskCfg, ok := a.cfg.Tasks[taskName]
	if !ok {
		return false, fmt.Sprintf("invalid task name %s, must be one of %v", taskName, a.cfg.TaskNames), nil
	}

	if task.Progress < 0 || len(taskCfg.Steps) < int(task.Progress) {
		return false, "invalid task progress", nil
	}

	for i := int(task.Progress); i < len(taskCfg.Steps); i++ {
		step := taskCfg.Steps[i]
		detail, err := a.processStep(step, taskCfg.WorkDir, params[1:])
		if err != nil {
			return false, fmt.Sprintf("step at index %d(%s) failed, %v, %s", i, step.Name, err, detail), nil
		}

		detail = fmt.Sprintf("step at index %d(%s) done:\n%s", i, step.Name, detail)
		if len(detail) > types.MaxTaskDetailLen {
			detail = detail[:types.MaxTaskDetailLen]
		}
		err = a.client.UpdateTask(task.ID, int8(i+1), detail)
		if err != nil {
			return false, "", fmt.Errorf("update task at step index %d failed, %v", i, err)
		}
	}

	return true, "", nil
}

func (a *Agent) processStep(step *config.Step, workDir string, params []string) (string, error) {
	var (
		out []byte
		err error
	)

	if step.Params {
		// For security, we don't use `bash -c`.
		out, err = goutils.ExeCmd(step.Cmd, params,
			func(c *exec.Cmd) {
				c.Dir = workDir
			})
	} else {
		// It is safe to use `bash -c`.
		out, err = goutils.ExeCmd("bash", []string{"-c", step.Cmd},
			func(c *exec.Cmd) {
				c.Dir = workDir
			})
	}

	return string(out), err
}

func (a *Agent) Destroy() {
	//
}

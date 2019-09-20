package config

import (
	"fmt"

	"github.com/tidyoux/goutils/viper"
)

type Step struct {
	Name   string
	Cmd    string
	Params bool
}

type Task struct {
	Name    string
	WorkDir string
	Steps   []*Step
}

type Config struct {
	URL       string
	WorkerID  uint64
	WorkerKey string
	Tasks     map[string]*Task
	TaskNames []string
}

func NewConfig() (*Config, error) {
	c := &Config{
		URL:       viper.GetString("url", "http://localhost"),
		WorkerID:  uint64(viper.GetInt64("workerID", 0)),
		WorkerKey: viper.GetString("workerKey", ""),
		Tasks:     make(map[string]*Task),
	}

	tasks, _ := viper.Get("tasks").([]interface{})
	if len(tasks) == 0 {
		return nil, fmt.Errorf("can't find tasks in config")
	}

	for i, task := range tasks {
		task, ok := task.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("parse task config at index %d failed, invalid format", i)
		}

		taskName, _ := task["name"].(string)
		if len(taskName) == 0 {
			return nil, fmt.Errorf("parse task config at index %d failed, name can't be empty", i)
		}

		if _, ok := c.Tasks[taskName]; ok {
			return nil, fmt.Errorf("parse task config at index %d failed, duplicate task name %s", i, taskName)
		}

		workDir, _ := task["workDir"].(string)

		steps, _ := task["steps"].([]interface{})
		if len(steps) == 0 {
			return nil, fmt.Errorf("parse task config at index %d failed, steps can't be empty", i)
		}

		t := &Task{
			Name:    taskName,
			WorkDir: workDir,
			Steps:   make([]*Step, 0, len(steps)),
		}

		for j, step := range steps {
			step, ok := step.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("parse task %s steps config at index %d failed, invalid format", taskName, j)
			}

			stepName, _ := step["name"].(string)
			if len(stepName) == 0 {
				return nil, fmt.Errorf("parse task %s steps config at index %d failed, name can't be empty", taskName, j)
			}

			cmd, _ := step["cmd"].(string)
			if len(cmd) == 0 {
				return nil, fmt.Errorf("parse task %s steps config at index %d failed, cmd can't be empty", taskName, j)
			}

			var params bool
			if p, ok := step["params"].(bool); ok {
				params = p
			}

			t.Steps = append(t.Steps, &Step{
				Name:   stepName,
				Cmd:    cmd,
				Params: params,
			})
		}

		c.Tasks[t.Name] = t
		c.TaskNames = append(c.TaskNames, t.Name)
	}

	return c, nil
}

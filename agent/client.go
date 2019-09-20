package agent

import (
	"strings"

	"github.com/tidyoux/router/common/client"
	"github.com/tidyoux/router/common/types"
)

type Client struct {
	*client.Client
}

func NewClient(url string) *Client {
	url = strings.TrimRight(url, "/") + "/v1/agent"
	return &Client{client.New(url)}
}

func (c *Client) Login(workerID uint64, workerKey string) error {
	var resp types.AgentLoginResponse
	err := c.Request("/login", &types.AgentLoginRequest{
		WorkerID:  workerID,
		WorkerKey: workerKey,
	}, &resp)
	if err != nil {
		return err
	}

	c.SetToken(resp.Token)
	return nil
}

func (c *Client) ListTask() (*types.ListTaskResponse, error) {
	var resp types.ListTaskResponse
	err := c.Request("/list-task", &types.AgentListTaskRequest{
		Token: c.Token(),
	}, &resp)
	return &resp, err
}

func (c *Client) AcceptTask(taskID uint64) error {
	return c.Request("/accept-task", &types.AcceptTaskRequest{
		Token:  c.Token(),
		TaskID: taskID,
	}, nil)
}

func (c *Client) UpdateTask(taskID uint64, progress int8, detail string) error {
	return c.Request("/update-task", &types.UpdateTaskRequest{
		Token:    c.Token(),
		TaskID:   taskID,
		Progress: progress,
		Detail:   detail,
	}, nil)
}

func (c *Client) FinishTask(taskID uint64, success bool, detail string) error {
	return c.Request("/finish-task", &types.FinishTaskRequest{
		Token:   c.Token(),
		TaskID:  taskID,
		Success: success,
		Detail:  detail,
	}, nil)
}

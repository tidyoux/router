package client

import (
	"strings"

	"github.com/tidyoux/router/common/client"
	"github.com/tidyoux/router/common/crypto"
	"github.com/tidyoux/router/common/types"
)

type Client struct {
	*client.Client
}

func New(url string) *Client {
	url = strings.TrimRight(url, "/") + "/v1/user"
	return &Client{client.New(url)}
}

func (c *Client) Login(username, password string) (uint64, error) {
	password = crypto.Sum([]byte(password)).String()

	var resp types.UserLoginResponse
	err := c.Request("/login", &types.UserLoginRequest{
		Username: username,
		Password: password,
	}, &resp)
	if err != nil {
		return 0, err
	}

	c.SetToken(resp.Token)
	return resp.UserID, nil
}

func (c *Client) Logout() error {
	err := c.Request("/logout", &types.UserLogoutRequest{
		Token: c.Token(),
	}, nil)
	if err != nil {
		return err
	}

	c.SetToken("")
	return nil
}

func (c *Client) Ping() error {
	return c.Request("/ping", &types.UserPingRequest{
		Token: c.Token(),
	}, nil)
}

func (c *Client) UpdatePassword(password string) error {
	password = crypto.Sum([]byte(password)).String()

	err := c.Request("/update-password", &types.UserUpdatePasswordRequest{
		Token:    c.Token(),
		Password: password,
	}, nil)
	if err != nil {
		return err
	}

	c.SetToken("")
	return nil
}

func (c *Client) ListUser() ([]*types.User, error) {
	var resp []*types.User
	err := c.Request("/list-user", &types.ListUserRequest{
		Token: c.Token(),
	}, &resp)
	return resp, err
}

func (c *Client) UpdateUserDetail(userID uint64, detail string) error {
	return c.Request("/update-user-detail", &types.UserUpdateUserRequest{
		Token:  c.Token(),
		UserID: userID,
		Detail: detail,
	}, nil)
}

func (c *Client) ResetUserPassword(userID uint64) error {
	return c.Request("/reset-user-password", &types.UserUpdateUserRequest{
		Token:  c.Token(),
		UserID: userID,
	}, nil)
}

func (c *Client) EnableUser(userID uint64) error {
	return c.Request("/enable-user", &types.UserUpdateUserRequest{
		Token:  c.Token(),
		UserID: userID,
	}, nil)
}

func (c *Client) DisableUser(userID uint64) error {
	return c.Request("/disable-user", &types.UserUpdateUserRequest{
		Token:  c.Token(),
		UserID: userID,
	}, nil)
}

func (c *Client) ListWorker() ([]*types.Worker, error) {
	var resp []*types.Worker
	err := c.Request("/list-worker", &types.ListWorkerRequest{
		Token: c.Token(),
	}, &resp)
	return resp, err
}

func (c *Client) AddWorker(name, desc string) (*types.AddWorkerResponse, error) {
	var resp types.AddWorkerResponse
	err := c.Request("/add-worker", &types.AddWorkerRequest{
		Token: c.Token(),
		Name:  name,
		Desc:  desc,
	}, &resp)
	return &resp, err
}

func (c *Client) EnableWorker(workerID uint64) error {
	return c.Request("/enable-worker", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
	}, nil)
}

func (c *Client) DisableWorker(workerID uint64) error {
	return c.Request("/disable-worker", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
	}, nil)
}

func (c *Client) UpdateWorkerName(workerID uint64, name string) error {
	return c.Request("/update-worker-name", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		Name:     name,
	}, nil)
}

func (c *Client) UpdateWorkerDesc(workerID uint64, desc string) error {
	return c.Request("/update-worker-desc", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		Desc:     desc,
	}, nil)
}

func (c *Client) AddWorkerUser(workerID, userID uint64) error {
	return c.Request("/add-worker-user", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		UserID:   userID,
	}, nil)
}

func (c *Client) RemoveWorkerUser(workerID, userID uint64) error {
	return c.Request("/remove-worker-user", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		UserID:   userID,
	}, nil)
}

func (c *Client) RemoveWorker(workerID uint64) error {
	return c.Request("/remove-worker", &types.UpdateWorkerRequest{
		Token:    c.Token(),
		WorkerID: workerID,
	}, nil)
}

func (c *Client) ListTask(workerID uint64, offset, limit int64) (*types.ListTaskResponse, error) {
	var resp types.ListTaskResponse
	err := c.Request("/list-task", &types.UserListTaskRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		Offset:   offset,
		Limit:    limit,
	}, &resp)
	return &resp, err
}

func (c *Client) SendTask(workerID uint64, params string) (*types.SendTaskResponse, error) {
	var resp types.SendTaskResponse
	err := c.Request("/send-task", &types.SendTaskRequest{
		Token:    c.Token(),
		WorkerID: workerID,
		Params:   params,
	}, &resp)
	return &resp, err
}

func (c *Client) TaskStatus(taskID uint64) (*types.TaskStatusResponse, error) {
	var resp types.TaskStatusResponse
	err := c.Request("/task-status", &types.TaskStatusRequest{
		Token:  c.Token(),
		TaskID: taskID,
	}, &resp)
	return &resp, err
}

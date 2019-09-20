package client_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidyoux/router/client"
)

var (
	c = client.New("http://localhost:8082")

	username = "admin"
	password = "123456"
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var p = fmt.Println
var pf = fmt.Printf

func pJSON(v interface{}) {
	buf, _ := json.Marshal(v)
	p(string(buf))
}

func login(t *testing.T) {
	_, err := c.Login(username, password)
	assert(t, err)
}

func TestPing(t *testing.T) {
	login(t)

	err := c.Ping()
	assert(t, err)
}

func TestUpdatePassword(t *testing.T) {
	login(t)

	err := c.UpdatePassword("123456")
	assert(t, err)
}

func TestListUser(t *testing.T) {
	login(t)

	users, err := c.ListUser()
	assert(t, err)

	p("total:", len(users))
	for _, u := range users {
		pJSON(u)
	}
}

func TestUpdateUserDetail(t *testing.T) {
	login(t)

	err := c.UpdateUserDetail(1, "a test user.")
	assert(t, err)
}

func TestResetUserPassword(t *testing.T) {
	login(t)

	err := c.ResetUserPassword(1)
	assert(t, err)
}

func TestEnableUser(t *testing.T) {
	login(t)

	err := c.EnableUser(1)
	assert(t, err)
}

func TestDisableUser(t *testing.T) {
	login(t)

	err := c.DisableUser(1)
	assert(t, err)
}

func TestListWorker(t *testing.T) {
	login(t)

	works, err := c.ListWorker()
	assert(t, err)

	p("total:", len(works))
	for _, w := range works {
		pJSON(w)
	}
}

func TestAddWorker(t *testing.T) {
	login(t)

	resp, err := c.AddWorker("test", "the test worker.")
	assert(t, err)

	pJSON(resp)
}

func TestEnableWorker(t *testing.T) {
	login(t)

	err := c.EnableWorker(1)
	assert(t, err)
}

func TestDisableWorker(t *testing.T) {
	login(t)

	err := c.DisableWorker(1)
	assert(t, err)
}

func TestUpdateWorker(t *testing.T) {
	login(t)

	err := c.UpdateWorkerName(2, "tester")
	assert(t, err)

	err = c.UpdateWorkerDesc(2, "test1's desc.")
	assert(t, err)
}

func TestAddWorkerUser(t *testing.T) {
	login(t)

	err := c.AddWorkerUser(1, 1)
	assert(t, err)
}

func TestRemoveWorkerUser(t *testing.T) {
	login(t)

	err := c.RemoveWorkerUser(1, 1)
	assert(t, err)
}

func TestRemoveWorker(t *testing.T) {
	login(t)

	err := c.RemoveWorker(6)
	assert(t, err)
}

func TestListTask(t *testing.T) {
	login(t)

	resp, err := c.ListTask(1, 0, 10)
	assert(t, err)

	p("total:", resp.Total)
	for _, task := range resp.Tasks {
		pJSON(task)
	}
}

func TestSendTask(t *testing.T) {
	login(t)

	resp, err := c.SendTask(3, "a:100,b:12,c:abc")
	assert(t, err)

	pJSON(resp)
}

func TestTaskStatus(t *testing.T) {
	login(t)

	resp, err := c.TaskStatus(1)
	assert(t, err)

	pJSON(resp)
}

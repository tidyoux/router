package cache

import "github.com/tidyoux/router/common/types"

var (
	defaultCache = &Cache{}
)

func C() *Cache {
	return defaultCache
}

type Cache struct {
	username string
	users    []*types.User
	userIdx  map[uint64]*types.User

	workers   []*types.Worker
	workerIdx map[uint64]*types.Worker

	taskWorkerID uint64
	totalTask    int64
	tasks        []*types.Task
}

func (c *Cache) Username() string {
	return c.username
}

func (c *Cache) SetUsername(username string) {
	c.username = username
}

func (c *Cache) Users() []*types.User {
	return c.users
}

func (c *Cache) UserByID(userID uint64) (*types.User, bool) {
	if len(c.userIdx) == 0 {
		return nil, false
	}

	user, ok := c.userIdx[userID]
	return user, ok
}

func (c *Cache) SetUsers(users []*types.User) {
	c.users = users
	c.userIdx = make(map[uint64]*types.User, len(users))
	for _, u := range users {
		c.userIdx[u.ID] = u
	}
}

func (c *Cache) Workers() []*types.Worker {
	return c.workers
}

func (c *Cache) WorkerByID(workerID uint64) (*types.Worker, bool) {
	if len(c.workerIdx) == 0 {
		return nil, false
	}

	worker, ok := c.workerIdx[workerID]
	return worker, ok
}

func (c *Cache) SetWorkers(workers []*types.Worker) {
	c.workers = workers
	c.workerIdx = make(map[uint64]*types.Worker, len(workers))
	for _, w := range workers {
		c.workerIdx[w.ID] = w
	}
}

func (c *Cache) TaskWorkerID() uint64 {
	return c.taskWorkerID
}

func (c *Cache) TotalTask() int64 {
	return c.totalTask
}

func (c *Cache) Tasks() []*types.Task {
	return c.tasks
}

func (c *Cache) SetTasks(workerID uint64, total int64, tasks []*types.Task) {
	c.taskWorkerID = workerID
	c.totalTask = total
	c.tasks = tasks
}

func (c *Cache) Clear() {
	c.username = ""
	c.users = nil
	c.userIdx = nil

	c.workers = nil
	c.workerIdx = nil

	c.taskWorkerID = 0
	c.totalTask = 0
	c.tasks = nil
}

package handler

import (
	"fmt"
	"strings"

	"github.com/tidyoux/router/common/crypto"
	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server"
	"github.com/tidyoux/router/server/model"
	"github.com/jinzhu/gorm"
)

const (
	MinUsernameLen     = 1
	MaxUsernameLen     = 32
	MinUserPasswordLen = 6

	MaxTaskParamsLen = 1024
)

var (
	errInvalidUserToken = fmt.Errorf("invalid user token")
	errUserDisabled     = fmt.Errorf("user disabled")

	errWorkerDisabled = fmt.Errorf("worker disabled")
	errWorkerOwner    = fmt.Errorf("you don't own this worker")
)

func init() {
	validUsernamePassword := func(name, password string) error {
		if len(name) < MinUsernameLen || MaxUsernameLen < len(name) {
			return fmt.Errorf("invalid user name length %d, should in [%d, %d]",
				len(name), MinUsernameLen, MaxUsernameLen)
		}

		if len(password) < MinUserPasswordLen {
			return fmt.Errorf("invalid user password length: %d, should >= %d",
				len(password), MinUserPasswordLen)
		}

		return nil
	}

	H("/user/login", func(req *types.UserLoginRequest) (*types.UserLoginResponse, error) {
		req.Username = strings.TrimSpace(req.Username)
		req.Password = strings.TrimSpace(req.Password)

		err := validUsernamePassword(req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		userID, token, err := server.S.UserLogin(req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return &types.UserLoginResponse{
			UserID: userID,
			Token:  token,
		}, nil
	})

	H("/user/logout", func(req *types.UserLogoutRequest) (bool, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return false, err
		}

		server.S.UserLogout(uint64(user.ID))
		return true, nil
	})

	H("/user/ping", func(req *types.UserPingRequest) (bool, error) {
		if _, err := validUser(req.Token); err != nil {
			return false, err
		}

		return true, nil
	})

	H("/user/update-password", func(req *types.UserUpdatePasswordRequest) (bool, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return false, err
		}

		err = validUsernamePassword(user.Name, req.Password)
		if err != nil {
			return false, err
		}

		password := crypto.Sum([]byte(req.Password)).String()
		if password != user.Password {
			err = user.UpdatePassword(password)
			if err != nil {
				return false, fmt.Errorf("db update user password failed, %v", err)
			}
		}

		server.S.UserLogout(uint64(user.ID))
		return true, nil
	})

	H("/user/list-worker", func(req *types.ListWorkerRequest) (types.ListWorkerResponse, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return nil, err
		}

		var workers []*model.Worker
		if user.IsAdmin() {
			workers, err = model.FindAllWorkers()
		} else {
			workers, err = model.FindWorkersByUserID(uint64(user.ID))
		}
		if err != nil {
			return nil, fmt.Errorf("db find workers failed, %v", err)
		}

		resp := make([]*types.Worker, 0, len(workers))
		for _, w := range workers {
			userIDs, err := model.FindWorkerUsers(uint64(w.ID))
			if err != nil {
				return nil, fmt.Errorf("db find worker users failed, %v", err)
			}

			resp = append(resp, &types.Worker{
				ID:        uint64(w.ID),
				Key:       w.Key,
				Name:      w.Name,
				Desc:      w.Desc,
				Status:    w.Status,
				CreatedAt: w.CreatedAt.Unix(),
				Users:     userIDs,
			})
		}
		return resp, nil
	})

	H("/user/update-worker-name", func(req *types.UpdateWorkerRequest) (bool, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return false, err
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		if err := validWorkerOwner(user, req.WorkerID); err != nil {
			return false, err
		}

		if err := validWorkerEnabled(worker.Status); err != nil {
			return false, err
		}

		req.Name = strings.TrimSpace(req.Name)

		if len(req.Name) == 0 || req.Name == worker.Name {
			return true, nil
		}

		err = validWorkerNameDesc(req.Name, worker.Desc)
		if err != nil {
			return false, err
		}

		if w, err := model.FindWorkerByName(req.Name); err != nil {
			if err != gorm.ErrRecordNotFound {
				return false, fmt.Errorf("db find worker by name failed, %v", err)
			}
		} else if w.ID != worker.ID {
			return false, fmt.Errorf("duplicate worker name")
		}

		err = worker.UpdateName(req.Name)
		if err != nil {
			return false, fmt.Errorf("db update worker name failed, %v", err)
		}

		return true, nil
	})

	H("/user/update-worker-desc", func(req *types.UpdateWorkerRequest) (bool, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return false, err
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		if err := validWorkerOwner(user, req.WorkerID); err != nil {
			return false, err
		}

		if err := validWorkerEnabled(worker.Status); err != nil {
			return false, err
		}

		req.Desc = strings.TrimSpace(req.Desc)

		if len(req.Desc) == 0 || req.Desc == worker.Desc {
			return true, nil
		}

		err = validWorkerNameDesc(worker.Name, req.Desc)
		if err != nil {
			return false, err
		}

		err = worker.UpdateDesc(req.Desc)
		if err != nil {
			return false, fmt.Errorf("db update worker desc failed, %v", err)
		}

		return true, nil
	})

	H("/user/list-task", func(req *types.UserListTaskRequest) (*types.ListTaskResponse, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return nil, err
		}

		if _, err := model.FindWorkerByID(req.WorkerID); err != nil {
			return nil, fmt.Errorf("db find worker failed, %v", err)
		}

		if err := validWorkerOwner(user, req.WorkerID); err != nil {
			return nil, err
		}

		total, err := model.FindTaskCountByWorkerID(req.WorkerID)
		if err != nil {
			return nil, fmt.Errorf("db find task count failed, %v", err)
		}

		tasks, err := model.FindTasksByWorkerID(req.WorkerID, req.Offset, req.Limit)
		if err != nil {
			return nil, fmt.Errorf("db find tasks failed, %v", err)
		}

		resp := &types.ListTaskResponse{
			Total: total,
			Tasks: make([]*types.Task, 0, len(tasks)),
		}

		for _, t := range tasks {
			user, _ := model.FindUserByID(t.UserID)
			resp.Tasks = append(resp.Tasks, &types.Task{
				ID:        uint64(t.ID),
				Params:    t.Params,
				Creator:   user.Name,
				Status:    t.Status,
				Progress:  t.Progress,
				Detail:    t.Detail,
				CreatedAt: t.CreatedAt.Unix(),
				UpdatedAt: t.UpdatedAt.Unix(),
			})
		}

		return resp, nil
	})

	H("/user/send-task", func(req *types.SendTaskRequest) (*types.SendTaskResponse, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return nil, err
		}

		req.Params = strings.TrimSpace(req.Params)

		if len(req.Params) > MaxTaskParamsLen {
			return nil, fmt.Errorf("invalid task params length: %d, should <= %d",
				len(req.Params), MaxTaskParamsLen)
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return nil, fmt.Errorf("db find worker failed, %v", err)
		}

		if err := validWorkerOwner(user, req.WorkerID); err != nil {
			return nil, err
		}

		if err := validWorkerEnabled(worker.Status); err != nil {
			return nil, err
		}

		task := model.NewTask(uint64(user.ID), req.WorkerID, req.Params)
		err = task.Insert()
		if err != nil {
			return nil, fmt.Errorf("db insert task failed, %v", err)
		}

		return &types.SendTaskResponse{
			TaskID: uint64(task.ID),
		}, nil
	})

	H("/user/task-status", func(req *types.TaskStatusRequest) (*types.TaskStatusResponse, error) {
		user, err := validUser(req.Token)
		if err != nil {
			return nil, err
		}

		task, err := model.FindTaskByID(req.TaskID)
		if err != nil {
			return nil, fmt.Errorf("db find task failed, %v", err)
		}

		if err := validWorkerOwner(user, task.WorkerID); err != nil {
			return nil, err
		}

		return &types.TaskStatusResponse{
			Status:   task.Status,
			Progress: task.Progress,
			Detail:   task.Detail,
		}, nil
	})
}

func validUser(token string) (*model.User, error) {
	userID, ok := server.S.ValidUserToken(token)
	if !ok {
		return nil, errInvalidUserToken
	}

	user, err := model.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("db find user failed, %v", err)
	}

	if user.Status == types.UserDisabled {
		return nil, errUserDisabled
	}

	return user, nil
}

func validWorkerEnabled(workerStatus int8) error {
	if workerStatus != types.WorkerEnabled {
		return errWorkerDisabled
	}

	return nil
}

func validWorkerOwner(user *model.User, workerID uint64) error {
	if user.IsAdmin() {
		return nil
	}

	if _, err := model.FindUserWorker(uint64(user.ID), workerID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errWorkerOwner
		}

		return fmt.Errorf("db find user worker failed, %v", err)
	}

	return nil
}

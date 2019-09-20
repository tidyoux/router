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
	MaxUserDetailLen = 1024

	MinWorkerNameLen = 1
	MaxWorkerNameLen = 32
	MaxWorkerDescLen = 1024
)

var (
	errPermissionDenied = fmt.Errorf("permission denied")
)

func init() {
	H("/user/list-user", func(req *types.ListUserRequest) (types.ListUserResponse, error) {
		if err := validAdmin(req.Token); err != nil {
			return nil, err
		}

		users, err := model.FindAllUsers()
		if err != nil {
			return nil, fmt.Errorf("db find users failed, %v", err)
		}

		resp := make([]*types.User, 0, len(users))
		for _, u := range users {
			if u.Name != types.AdminUsername {
				resp = append(resp, &types.User{
					ID:     uint64(u.ID),
					Name:   u.Name,
					Detail: u.Detail,
					Status: u.Status,
				})
			}
		}
		return resp, nil
	})

	H("user/update-user-detail", func(req *types.UserUpdateUserRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		req.Detail = strings.TrimSpace(req.Detail)
		if len(req.Detail) > MaxUserDetailLen {
			return false, fmt.Errorf("invalid user name length %d, should <= %d",
				len(req.Detail), MaxUserDetailLen)
		}

		user, err := model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		if req.Detail != user.Detail {
			err = user.UpdateDetail(req.Detail)
			if err != nil {
				return false, fmt.Errorf("db update user detail failed, %v", err)
			}
		}

		return true, nil
	})

	H("/user/reset-user-password", func(req *types.UserUpdateUserRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		user, err := model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		password := types.DefaultPassword
		password = crypto.Sum([]byte(password)).String()
		if password != user.Password {
			err = user.UpdatePassword(password)
			if err != nil {
				return false, fmt.Errorf("db update user password failed, %v", err)
			}
		}

		server.S.UserLogout(req.UserID)
		return true, nil
	})

	H("/user/enable-user", func(req *types.UserUpdateUserRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		user, err := model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		if user.Status == types.UserEnabled {
			return true, nil
		}

		err = user.Enable()
		if err != nil {
			return false, fmt.Errorf("db enable user failed, %v", err)
		}

		return true, nil
	})

	H("/user/disable-user", func(req *types.UserUpdateUserRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		user, err := model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		if user.Status == types.UserDisabled {
			return true, nil
		}

		err = user.Disable()
		if err != nil {
			return false, fmt.Errorf("db disable user failed, %v", err)
		}

		return true, nil
	})

	H("/user/add-worker", func(req *types.AddWorkerRequest) (*types.AddWorkerResponse, error) {
		if err := validAdmin(req.Token); err != nil {
			return nil, err
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Desc = strings.TrimSpace(req.Desc)

		err := validWorkerNameDesc(req.Name, req.Desc)
		if err != nil {
			return nil, err
		}

		if worker, err := model.FindWorkerByName(req.Name); err == nil {
			if req.Desc != worker.Desc {
				return nil, fmt.Errorf("worker %s already existed", req.Name)
			}

			return &types.AddWorkerResponse{
				WorkerID:  uint64(worker.ID),
				WorkerKey: worker.Key,
			}, nil
		}

		key := server.GenToken()
		worker := model.NewWorker(key, req.Name, req.Desc)
		err = worker.Insert()
		if err != nil {
			return nil, fmt.Errorf("db insert worker failed, %v", err)
		}

		return &types.AddWorkerResponse{
			WorkerID:  uint64(worker.ID),
			WorkerKey: key,
		}, nil
	})

	H("/user/enable-worker", func(req *types.UpdateWorkerRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		if worker.Status == types.WorkerEnabled {
			return true, nil
		}

		err = worker.Enable()
		if err != nil {
			return false, fmt.Errorf("db enable worker failed, %v", err)
		}

		return true, nil
	})

	H("/user/disable-worker", func(req *types.UpdateWorkerRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		if worker.Status == types.WorkerDisabled {
			return true, nil
		}

		err = worker.Disable()
		if err != nil {
			return false, fmt.Errorf("db disable worker failed, %v", err)
		}

		return true, nil
	})

	H("/user/add-worker-user", func(req *types.UpdateWorkerRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		_, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		_, err = model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		_, err = model.FindUserWorker(req.UserID, req.WorkerID)
		if err == nil {
			return true, nil
		}

		if err != gorm.ErrRecordNotFound {
			return false, fmt.Errorf("db find user worker failed, %v", err)
		}

		userWorker := model.NewUserWorker(req.UserID, req.WorkerID)
		err = userWorker.Insert()
		if err != nil {
			return false, fmt.Errorf("db insert user worker failed, %v", err)
		}

		return true, nil
	})

	H("/user/remove-worker-user", func(req *types.UpdateWorkerRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		_, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		_, err = model.FindUserByID(req.UserID)
		if err != nil {
			return false, fmt.Errorf("db find user failed, %v", err)
		}

		userWorker, err := model.FindUserWorker(req.UserID, req.WorkerID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return true, nil
			}

			return false, fmt.Errorf("db find user worker failed, %v", err)
		}

		err = userWorker.Delete()
		if err != nil {
			return false, fmt.Errorf("db delete user worker failed, %v", err)
		}

		return true, nil
	})

	H("/user/remove-worker", func(req *types.UpdateWorkerRequest) (bool, error) {
		if err := validAdmin(req.Token); err != nil {
			return false, err
		}

		worker, err := model.FindWorkerByID(req.WorkerID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return true, nil
			}

			return false, fmt.Errorf("db find worker failed, %v", err)
		}

		err = model.DeleteWorkerUsers(req.WorkerID)
		if err != nil {
			return false, fmt.Errorf("db delete worker users failed, %v", err)
		}

		err = worker.Delete()
		if err != nil {
			return false, fmt.Errorf("db delete worker failed, %v", err)
		}

		server.S.RemoveWorker(req.WorkerID)

		return true, nil
	})
}

func validAdmin(token string) error {
	user, err := validUser(token)
	if err != nil {
		return err
	}

	if !user.IsAdmin() {
		return errPermissionDenied
	}

	return nil
}

func validWorkerNameDesc(name, desc string) error {
	if len(name) < MinWorkerNameLen || MaxWorkerNameLen < len(name) {
		return fmt.Errorf("invalid worker name length %d, should in [%d, %d]",
			len(name), MinWorkerNameLen, MaxWorkerNameLen)
	}

	if len(desc) > MaxWorkerDescLen {
		return fmt.Errorf("invalid worker desc length: %d, should <= %d",
			len(desc), MaxWorkerDescLen)
	}

	return nil
}

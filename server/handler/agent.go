package handler

import (
	"fmt"
	"strings"

	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server"
	"github.com/tidyoux/router/server/model"
)

var (
	ErrInvalidAgentToken = fmt.Errorf("invalid agent token")
)

func init() {
	H("/agent/login", func(req *types.AgentLoginRequest) (*types.AgentLoginResponse, error) {
		token, err := server.S.AgentLogin(req.WorkerID, req.WorkerKey)
		if err != nil {
			return nil, err
		}

		return &types.AgentLoginResponse{token}, nil
	})

	validAgent := func(token string) (uint64, error) {
		workerID, ok := server.S.ValidAgentToken(token)
		if !ok {
			return 0, ErrInvalidAgentToken
		}

		worker, err := model.FindWorkerByID(workerID)
		if err != nil {
			return 0, fmt.Errorf("db find worker failed, %v", err)
		}

		if err := validWorkerEnabled(worker.Status); err != nil {
			return 0, err
		}

		return workerID, nil
	}
	H("/agent/list-task", func(req *types.AgentListTaskRequest) (*types.ListTaskResponse, error) {
		workerID, err := validAgent(req.Token)
		if err != nil {
			return nil, err
		}

		tasks, err := model.FindAllUnfinishedTasksByWorkerID(workerID)
		if err != nil {
			return nil, fmt.Errorf("db find tasks failed, %v", err)
		}

		resp := &types.ListTaskResponse{
			Total: int64(len(tasks)),
			Tasks: make([]*types.Task, 0, len(tasks)),
		}

		for _, t := range tasks {
			resp.Tasks = append(resp.Tasks, &types.Task{
				ID:        uint64(t.ID),
				Params:    t.Params,
				Status:    t.Status,
				Progress:  t.Progress,
				CreatedAt: t.CreatedAt.Unix(),
			})
		}

		return resp, nil
	})

	H("/agent/accept-task", func(req *types.AcceptTaskRequest) (bool, error) {
		if _, err := validAgent(req.Token); err != nil {
			return false, err
		}

		task, err := model.FindTaskByID(req.TaskID)
		if err != nil {
			return false, fmt.Errorf("db find task failed, %v", err)
		}

		switch task.Status {
		case types.TaskAccepted:
			return true, nil
		case types.TaskFailed, types.TaskSuccess:
			return false, fmt.Errorf("can't accept a finished task")
		default:
			err = task.Accept()
			if err != nil {
				return false, fmt.Errorf("db accept task failed, %v", err)
			}

			return true, nil
		}
	})

	H("/agent/update-task", func(req *types.UpdateTaskRequest) (bool, error) {
		if _, err := validAgent(req.Token); err != nil {
			return false, err
		}

		req.Detail = strings.TrimSpace(req.Detail)

		if len(req.Detail) > types.MaxTaskDetailLen {
			return false, fmt.Errorf("invalid task detail length: %d, should <= %d",
				len(req.Detail), types.MaxTaskDetailLen)
		}

		task, err := model.FindTaskByID(req.TaskID)
		if err != nil {
			return false, fmt.Errorf("db find task failed, %v", err)
		}

		switch task.Status {
		case types.TaskRecord:
			return false, fmt.Errorf("task must accept first")
		case types.TaskFailed, types.TaskSuccess:
			return false, fmt.Errorf("can't update a finished task")
		default:
			err = task.UpdateProgress(req.Progress, req.Detail)
			if err != nil {
				return false, fmt.Errorf("db update task progress failed, %v", err)
			}

			return true, nil
		}
	})

	H("/agent/finish-task", func(req *types.FinishTaskRequest) (bool, error) {
		if _, err := validAgent(req.Token); err != nil {
			return false, err
		}

		req.Detail = strings.TrimSpace(req.Detail)

		if len(req.Detail) > types.MaxTaskDetailLen {
			return false, fmt.Errorf("invalid task detail length: %d, should <= %d",
				len(req.Detail), types.MaxTaskDetailLen)
		}

		task, err := model.FindTaskByID(req.TaskID)
		if err != nil {
			return false, fmt.Errorf("db find task failed, %v", err)
		}

		switch task.Status {
		case types.TaskRecord:
			return false, fmt.Errorf("task must accept first")
		case types.TaskSuccess:
			if !req.Success || req.Detail != task.Detail {
				return false, fmt.Errorf("task has finished already")
			}

			return true, nil
		case types.TaskFailed:
			if req.Success || req.Detail != task.Detail {
				return false, fmt.Errorf("task has finished already")
			}

			return true, nil
		default:
			if req.Success {
				err = task.Success(req.Detail)
			} else {
				err = task.Fail(req.Detail)
			}
			if err != nil {
				return false, fmt.Errorf("db finish task failed, %v", err)
			}

			return true, nil
		}
	})
}

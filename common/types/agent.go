package types

type AgentLoginRequest struct {
	WorkerID  uint64 `json:"worker_id"`
	WorkerKey string `json:"worker_key"`
}

type AgentLoginResponse struct {
	Token string `json:"token"`
}

type AgentListTaskRequest struct {
	Token string `json:"token"`
}

type AcceptTaskRequest struct {
	Token  string `json:"token"`
	TaskID uint64 `json:"task_id"`
}

type UpdateTaskRequest struct {
	Token    string `json:"token"`
	TaskID   uint64 `json:"task_id"`
	Progress int8   `json:"progress"`
	Detail   string `json:"detail"`
}

type FinishTaskRequest struct {
	Token   string `json:"token"`
	TaskID  uint64 `json:"task_id"`
	Success bool   `json:"success"`
	Detail  string `json:"detail"`
}

package types

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

type UserLogoutRequest struct {
	Token string `json:"token"`
}

type UserPingRequest struct {
	Token string `json:"token"`
}

type UserUpdatePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ListUserRequest struct {
	Token string `json:"token"`
}

type User struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
	Status int8   `json:"status"`
}

type ListUserResponse []*User

type UserUpdateUserRequest struct {
	Token  string `json:"token"`
	UserID uint64 `json:"user_id"`
	Detail string `json:"detail"`
}

type ListWorkerRequest struct {
	Token string `json:"token"`
}

type Worker struct {
	ID        uint64 `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Status    int8   `json:"status"`
	CreatedAt int64  `json:"created_at"`

	Users []uint64 `json:"users"`
}
type ListWorkerResponse []*Worker

type AddWorkerRequest struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
}

type AddWorkerResponse struct {
	WorkerID  uint64 `json:"worker_id"`
	WorkerKey string `json:"worker_key"`
}

type UpdateWorkerRequest struct {
	Token    string `json:"token"`
	WorkerID uint64 `json:"worker_id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	UserID   uint64 `json:"user_id"`
}

type UserListTaskRequest struct {
	Token    string `json:"token"`
	WorkerID uint64 `json:"worker_id"`
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`
}

type Task struct {
	ID        uint64 `json:"id"`
	Params    string `json:"params"`
	Creator   string `json:"creator"`
	Status    int8   `json:"status"`
	Progress  int8   `json:"progress"`
	Detail    string `json:"detail"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ListTaskResponse struct {
	Total int64   `json:"total"`
	Tasks []*Task `json:"tasks"`
}

type SendTaskRequest struct {
	Token    string `json:"token"`
	WorkerID uint64 `json:"worker_id"`
	Params   string `json:"params"`
}

type SendTaskResponse struct {
	TaskID uint64 `json:"task_id"`
}

type TaskStatusRequest struct {
	Token  string `json:"token"`
	TaskID uint64 `json:"task_id"`
}

type TaskStatusResponse struct {
	Status   int8   `json:"status"`
	Progress int8   `json:"progress"`
	Detail   string `json:"detail"`
}

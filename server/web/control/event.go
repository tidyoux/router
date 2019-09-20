package control

// Event names.
const (
	ERerender = "rerender"

	ELogin  = "login"
	ELogout = "logout"

	EUpdatePassword = "update-password"

	EListWorker   = "list-worker"
	EUpdateWorker = "update-worker"

	EListTask = "list-task"
	ESendTask = "send-task"

	EListUser          = "list-user"
	EUpdateUserDetail  = "update-user-detail"
	EResetUserPassword = "reset-user-password"
	EEnableUser        = "enable-user"
	EDisableUser       = "disable-user"
	EAddWorker         = "add-worker"
	EEnableWorker      = "enable-worker"
	EDisableWorker     = "disable-worker"
	EAddWorkerUser     = "add-worker-user"
	ERemoveWorkerUser  = "remove-worker-user"
	ERemoveWorker      = "remove-worker"
)

type Event struct {
	name string
	data map[string]interface{}
}

func NewEvent(name string) *Event {
	return &Event{
		name: name,
		data: make(map[string]interface{}),
	}
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Set(k string, v interface{}) *Event {
	e.data[k] = v
	return e
}

func (e *Event) Get(k string) (interface{}, bool) {
	v, ok := e.data[k]
	return v, ok
}

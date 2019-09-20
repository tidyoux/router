package types

import "github.com/tidyoux/router/common/crypto"

// User status.
const (
	UserDisabled = 0
	UserEnabled  = 1
)

// Worker status.
const (
	WorkerDisabled = 0
	WorkerEnabled  = 1
)

// Task status.
const (
	TaskRecord   = 0
	TaskAccepted = 1
	TaskFailed   = 20
	TaskSuccess  = 50
)

// Content length.
const (
	MaxTaskDetailLen = 32 * 1024
)

const (
	AdminUsername = "admin"
)

var (
	DefaultPassword = crypto.Sum([]byte("123456")).String()
	InitialUserList = map[string]string{
		AdminUsername: DefaultPassword,
		"chris":       crypto.Sum(crypto.RandBytes(32)).String(),
		"kendall":     crypto.Sum(crypto.RandBytes(32)).String(),
		"arrow":       crypto.Sum(crypto.RandBytes(32)).String(),
		"claire":      crypto.Sum(crypto.RandBytes(32)).String(),
		"kerry":       crypto.Sum(crypto.RandBytes(32)).String(),
		"gene":        crypto.Sum(crypto.RandBytes(32)).String(),
		"joey":        crypto.Sum(crypto.RandBytes(32)).String(),
		"devin":       crypto.Sum(crypto.RandBytes(32)).String(),
		"ricki":       crypto.Sum(crypto.RandBytes(32)).String(),
	}
)

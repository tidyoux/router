package server

import (
	"fmt"
	"sync"

	"github.com/tidyoux/router/common/crypto"
	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/config"
	"github.com/tidyoux/router/server/model"
	"github.com/jinzhu/gorm"
)

var (
	errUsernameAlreadyExist = fmt.Errorf("username already exist")
	errUsernameNotExist     = fmt.Errorf("username not exist")
	errIncorrectPassword    = fmt.Errorf("incorrect password")
	errUserDisabled         = fmt.Errorf("user disabled")

	errIncorrectWorkerKey = fmt.Errorf("incorrect worker key")
	errWorkerDisabled     = fmt.Errorf("worker disabled")
)

var S *Server

func Init(cfg *config.Config) error {
	S = New(cfg)

	for username, password := range types.InitialUserList {
		_, _, err := S.UserRegister(username, password)
		if err != nil && err != errUsernameAlreadyExist {
			return fmt.Errorf("register user %s failed, %v", username, err)
		}
	}

	user, err := model.FindUserByName(types.AdminUsername)
	if err != nil {
		return fmt.Errorf("db find admin user failed, %v", err)
	}

	if user.Status == types.UserEnabled {
		return nil
	}

	err = user.Enable()
	if err != nil {
		return fmt.Errorf("db enable admin user failed, %v", err)
	}

	return nil
}

type Server struct {
	rwMutex sync.RWMutex

	cfg *config.Config

	users      map[string]uint64
	userTokens map[uint64]string

	agents      map[string]uint64
	agentTokens map[uint64]string
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,

		users:      make(map[string]uint64),
		userTokens: make(map[uint64]string),

		agents:      make(map[string]uint64),
		agentTokens: make(map[uint64]string),
	}
}

func (s *Server) UserRegister(username, password string) (uint64, string, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	password = crypto.Sum([]byte(password)).String()

	user, err := model.FindUserByName(username)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, "", fmt.Errorf("db find user failed, %v", err)
		}

		user = model.NewUser(username, password)
		err = user.Insert()
		if err != nil {
			return 0, "", fmt.Errorf("db insert user failed, %v", err)
		}
	} else {
		if password != user.Password {
			return 0, "", errUsernameAlreadyExist
		}
	}

	return uint64(user.ID), s.userLogin(uint64(user.ID)), nil
}

func (s *Server) UserLogin(username, password string) (uint64, string, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	user, err := model.FindUserByName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, "", errUsernameNotExist
		}

		return 0, "", fmt.Errorf("db find user failed, %v", err)
	}

	password = crypto.Sum([]byte(password)).String()
	if password != user.Password {
		return 0, "", errIncorrectPassword
	}

	if user.Status == types.UserDisabled {
		return 0, "", errUserDisabled
	}

	return uint64(user.ID), s.userLogin(uint64(user.ID)), nil
}

func (s *Server) userLogin(userID uint64) string {
	if token, ok := s.userTokens[userID]; ok {
		return token
	}

	token := GenToken()
	s.users[token] = userID
	s.userTokens[userID] = token
	return token
}

func (s *Server) UserLogout(userID uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if token, ok := s.userTokens[userID]; ok {
		delete(s.users, token)
	}
	delete(s.userTokens, userID)
}

func (s *Server) ValidUserToken(token string) (uint64, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	userID, ok := s.users[token]
	return userID, ok
}

func (s *Server) AgentLogin(workerID uint64, workerKey string) (string, error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if token, ok := s.agentTokens[workerID]; ok {
		return token, nil
	}

	worker, err := model.FindWorkerByID(workerID)
	if err != nil {
		return "", fmt.Errorf("db find worker failed, %v", err)
	}

	if workerKey != worker.Key {
		return "", errIncorrectWorkerKey
	}

	if worker.Status == types.WorkerDisabled {
		return "", errWorkerDisabled
	}

	token := GenToken()
	s.agents[token] = workerID
	s.agentTokens[workerID] = token
	return token, nil
}

func (s *Server) ValidAgentToken(token string) (uint64, bool) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	workerID, ok := s.agents[token]
	return workerID, ok
}

func (s *Server) RemoveWorker(workerID uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if token, ok := s.agentTokens[workerID]; ok {
		delete(s.agents, token)
	}
	delete(s.agentTokens, workerID)
}

package model

import (
	"sort"

	"github.com/tidyoux/router/common/db"
)

type UserWorker struct {
	Model

	UserID   uint64 `gorm:"index;unique_index:idx_user_worker"`
	WorkerID uint64 `gorm:"index;unique_index:idx_user_worker"`
}

func NewUserWorker(userID, workerID uint64) *UserWorker {
	return &UserWorker{
		UserID:   userID,
		WorkerID: workerID,
	}
}

func (*UserWorker) TableName() string { return "user_worker" }

func (u *UserWorker) Insert() error {
	return db.Default().Create(u).Error
}

func (u *UserWorker) Delete() error {
	return db.Default().Delete(u).Error
}

func FindUserWorker(userID, workerID uint64) (*UserWorker, error) {
	var userWorker UserWorker
	err := db.Default().First(&userWorker, "user_id = ? and worker_id = ?", userID, workerID).Error
	if err != nil {
		return nil, err
	}

	return &userWorker, nil
}

func FindUserWorkers(userID uint64) ([]uint64, error) {
	var workers []*UserWorker
	err := db.Default().Find(&workers, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	if len(workers) == 0 {
		return nil, nil
	}

	workerIDs := make([]uint64, 0, len(workers))
	for _, w := range workers {
		workerIDs = append(workerIDs, w.WorkerID)
	}

	sort.Slice(workerIDs, func(i, j int) bool {
		return workerIDs[i] < workerIDs[j]
	})

	return workerIDs, nil
}

func FindWorkerUsers(workerID uint64) ([]uint64, error) {
	var users []*UserWorker
	err := db.Default().Find(&users, "worker_id = ?", workerID).Error
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	userIDs := make([]uint64, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}

	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})

	return userIDs, nil
}

func DeleteWorkerUsers(workerID uint64) error {
	return db.Default().Where("worker_id = ?", workerID).Delete(UserWorker{}).Error
}

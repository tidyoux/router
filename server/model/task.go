package model

import (
	"github.com/tidyoux/router/common/db"
	"github.com/tidyoux/router/common/types"
)

type Task struct {
	Model

	UserID   uint64 `gorm:"index"`
	WorkerID uint64 `gorm:"index"`
	Params   string `gorm:"type:text"`
	Status   int8   `gorm:"type:tinyint;index"`
	Progress int8   `gorm:"type:tinyint;index"`
	Detail   string `gorm:"type:text"`
}

func NewTask(userID, workerID uint64, params string) *Task {
	return &Task{
		UserID:   userID,
		WorkerID: workerID,
		Params:   params,
	}
}

func (*Task) TableName() string { return "task" }

func (t *Task) Insert() error {
	return db.Default().Create(t).Error
}

func (t *Task) Accept() error {
	return t.update(M{
		"status": types.TaskAccepted,
	})
}

func (t *Task) UpdateProgress(progress int8, detail string) error {
	if progress < 0 {
		progress = 0
	}

	return t.update(M{
		"progress": progress,
		"detail":   detail,
	})
}

func (t *Task) Success(detail string) error {
	data := M{
		"status": types.TaskSuccess,
	}
	if len(detail) > 0 {
		data["detail"] = detail
	}
	return t.update(data)
}

func (t *Task) Fail(detail string) error {
	return t.update(M{
		"status": types.TaskFailed,
		"detail": detail,
	})
}

func (t *Task) update(values M) error {
	return db.Default().Model(t).Updates(values).Error
}

func FindTaskByID(id uint64) (*Task, error) {
	var task Task
	err := db.Default().First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func FindTaskCountByWorkerID(workerID uint64) (int64, error) {
	var count int64
	err := db.Default().Model(Task{}).Where("worker_id = ?", workerID).Count(&count).Error
	return count, err
}

func FindTasksByWorkerID(workerID uint64, offset, limit int64) ([]*Task, error) {
	var tasks []*Task
	err := db.Default().Offset(offset).Limit(limit).Order("id desc").Find(&tasks, "worker_id = ?", workerID).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func FindAllUnfinishedTasksByWorkerID(workerID uint64) ([]*Task, error) {
	var tasks []*Task
	err := db.Default().Find(&tasks, "worker_id = ? and status <= ?", workerID, types.TaskAccepted).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

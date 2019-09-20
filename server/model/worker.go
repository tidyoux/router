package model

import (
	"github.com/tidyoux/router/common/db"
	"github.com/tidyoux/router/common/types"
)

type Worker struct {
	Model

	Key    string `gorm:"size:64"`
	Name   string `gorm:"size:32;unique_index"`
	Desc   string `gorm:"type:text"`
	Status int8   `gorm:"type:tinyint;index"`
}

func NewWorker(key, name, desc string) *Worker {
	return &Worker{
		Key:    key,
		Name:   name,
		Desc:   desc,
		Status: types.WorkerDisabled,
	}
}

func (*Worker) TableName() string { return "worker" }

func (w *Worker) Insert() error {
	return db.Default().Create(w).Error
}

func (w *Worker) Enable() error {
	return w.update(M{
		"status": types.WorkerEnabled,
	})
}

func (w *Worker) Disable() error {
	return w.update(M{
		"status": types.WorkerDisabled,
	})
}

func (w *Worker) UpdateName(name string) error {
	return w.update(M{
		"name": name,
	})
}

func (w *Worker) UpdateDesc(desc string) error {
	return w.update(M{
		"desc": desc,
	})
}

func (w *Worker) Delete() error {
	return db.Default().Delete(w).Error
}

func (w *Worker) update(values M) error {
	return db.Default().Model(w).Updates(values).Error
}

func FindAllWorkers() ([]*Worker, error) {
	var workers []*Worker
	err := db.Default().Find(&workers).Error
	if err != nil {
		return nil, err
	}

	return workers, nil
}

func FindWorkersByUserID(userID uint64) ([]*Worker, error) {
	workerIDs, err := FindUserWorkers(userID)
	if err != nil {
		return nil, err
	}

	var workers []*Worker
	err = db.Default().Find(&workers, "id in (?)", workerIDs).Error
	if err != nil {
		return nil, err
	}

	return workers, nil
}

func FindWorkerByID(id uint64) (*Worker, error) {
	var worker Worker
	err := db.Default().First(&worker, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &worker, nil
}

func FindWorkerByName(name string) (*Worker, error) {
	var worker Worker
	err := db.Default().First(&worker, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &worker, nil
}

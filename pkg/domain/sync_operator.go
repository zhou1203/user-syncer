package domain

import (
	"context"

	"k8s.io/klog/v2"
)

type SyncerOperator struct {
	tasks []*Task
}

type Task struct {
	Syncer   Syncer
	Provider Provider
}

func NewSyncerOperator(tasks ...*Task) *SyncerOperator {
	return &SyncerOperator{tasks: tasks}
}

func (o *SyncerOperator) Sync(ctx context.Context) error {
	for _, t := range o.tasks {
		err := sync(ctx, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func sync(ctx context.Context, task *Task) error {
	objs, err := task.Provider.List(ctx)
	if err != nil {
		klog.Error(err)
		return err
	}

	for _, obj := range objs {
		err := task.Syncer.Sync(ctx, obj)
		if err != nil {
			klog.Error(err)
			continue
		}
	}
	return nil
}

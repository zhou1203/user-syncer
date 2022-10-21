package domain

import (
	"context"
	"k8s.io/klog/v2"
)

type SyncerOperator struct {
	Syncer Syncer
}

func NewSyncerOperator(syncer Syncer) *SyncerOperator {
	return &SyncerOperator{Syncer: syncer}
}

func (o *SyncerOperator) Sync(ctx context.Context, provider Provider) error {
	userList, err := provider.List(ctx)
	if err != nil {
		klog.Error(err)
		return err
	}
	for _, user := range userList {
		err := o.Syncer.Sync(ctx, user)
		if err != nil {
			klog.Error(err)
			continue
		}
	}

	return nil
}

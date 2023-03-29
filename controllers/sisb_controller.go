package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type SISBController struct {
	ctrl    controller.Controller
	stop    context.CancelFunc
	started bool
}

func NewSISBController(manager manager.Manager, reconciler reconcile.Reconciler) (*SISBController, error) {
	c, err := controller.NewUnmanaged("sisb-controller", manager, controller.Options{
		Reconciler:       reconciler,
		CacheSyncTimeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var siUnstructured, sbUnstructured unstructured.Unstructured
	siUnstructured.SetGroupVersionKind(instanceGvk)
	sbUnstructured.SetGroupVersionKind(bindingGvk)

	if err := c.Watch(&source.Kind{Type: &siUnstructured}, &handler.EnqueueRequestForObject{}); err != nil {
		return nil, err
	}

	if err := c.Watch(&source.Kind{Type: &sbUnstructured}, &handler.EnqueueRequestForObject{}); err != nil {
		return nil, err
	}

	return &SISBController{ctrl: c}, nil
}

func (c *SISBController) Start(ctx context.Context) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	c.stop = cancel
	go c.ctrl.Start(ctxWithCancel)
	c.started = true
}

func (c *SISBController) Stop(ctx context.Context) {
	logger := log.FromContext(ctx)
	logger.Info("--- STOP SISBController ---")
	c.stop()
	c.started = false
}

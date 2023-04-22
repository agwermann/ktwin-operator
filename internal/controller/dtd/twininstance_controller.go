/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dtd

import (
	"context"
	"fmt"

	twinservice "ktwin/operator/internal/resources/service"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dtdv0 "ktwin/operator/api/dtd/v0"
)

// TwinInstanceReconciler reconciles a TwinInstance object
type TwinInstanceReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	TwinService twinservice.TwinService
}

//+kubebuilder:rbac:groups=dtd.ktwin,resources=twininstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dtd.ktwin,resources=twininstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dtd.ktwin,resources=twininstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TwinInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *TwinInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	twinInstance := &dtdv0.TwinInstance{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, twinInstance)

	// Delete scenario
	if err != nil {
		if errors.IsNotFound(err) {
			return r.deleteTwinInstance(ctx, req, req.NamespacedName)
		}
		logger.Error(err, fmt.Sprintf("Unexpected error while deleting TwinInstance %s", req.Name))
		return ctrl.Result{}, err
	}

	return r.createUpdateTwinInstance(ctx, req, twinInstance)
}

func (r *TwinInstanceReconciler) createUpdateTwinInstance(ctx context.Context, req ctrl.Request, twinInstance *dtdv0.TwinInstance) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Create Service Instance
	kService := r.TwinService.GetService(twinInstance)
	err := r.Create(ctx, kService, &client.CreateOptions{})

	if err != nil {
		logger.Error(err, fmt.Sprintf("Error while creating Twin Instance %s", twinInstance.Spec.Id))
		return ctrl.Result{}, err
	}

	// Create MQTT Integrators

	// Create Triggers

	return ctrl.Result{}, nil
}

func (r *TwinInstanceReconciler) deleteTwinInstance(ctx context.Context, req ctrl.Request, namespacedName types.NamespacedName) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Create Service Instance
	kService := r.TwinService.GetDeletionService(namespacedName)
	err := r.Delete(ctx, kService, &client.DeleteOptions{})

	if err != nil {
		logger.Error(err, fmt.Sprintf("Error while deleting TwinInstance %s", namespacedName.Name))
		return ctrl.Result{}, err
	}

	// Create MQTT Integrators

	// Create Triggers

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TwinInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dtdv0.TwinInstance{}).
		Complete(r)
}

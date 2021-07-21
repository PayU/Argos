/*
Copyright 2021.

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

package controllers

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	argosv1 "slo.payu.com/argos/api/v1"
	renderer "slo.payu.com/argos/renderer"
)

// SloReconciler reconciles a Slo object
type SloReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argos.slo.payu.com,resources=sloes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argos.slo.payu.com,resources=sloes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argos.slo.payu.com,resources=sloes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Slo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SloReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var slo argosv1.Slo
	var err error
	var log = log.FromContext(ctx)

	if err = r.Get(context.Background(), req.NamespacedName, &slo); err != nil {
		log.Error(err, "SLO Controller faced a problem when fetching desired SLO")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	sloJSON, err := json.Marshal(slo)

	if err != nil {
		log.Error(err, "SLO controller faced problem marshaling SLO struct")
	}

	log.Info(string(sloJSON))
	renderer.RenderTemplate()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argosv1.Slo{}).
		Complete(r)
}

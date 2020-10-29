/*

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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

// BareMetalHostReconciler reconciles a BareMetalHost object
type BareMetalHostReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *BareMetalHostReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("baremetalhost", req.NamespacedName)

	logger.Info("changed")

	return ctrl.Result{}, nil
}

func (r *BareMetalHostReconciler) SetupWithManager(mgr ctrl.Manager) error {

	mapper := hostMapper{
		client: mgr.GetClient(),
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&bmh.BareMetalHost{}).
		Named("baremetalhost").
		Watches(&source.Kind{Type: &hwcc.HardwareClassification{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: &mapper}).
		Complete(r)
}

type hostMapper struct {
	client client.Client
}

func (m *hostMapper) Map(obj handler.MapObject) []ctrl.Request {
	log := ctrl.Log.WithName("controllers").WithName("BareMetalHost").WithName("mapper").
		WithValues("HardwareClassification",
			fmt.Sprintf("%s/%s", obj.Meta.GetNamespace(), obj.Meta.GetName()))

	bmhHostList := bmh.BareMetalHostList{}
	opts := &client.ListOptions{}
	err := m.client.List(context.TODO(), &bmhHostList, opts)
	if err != nil {
		log.Error(err, "could not fetch host list")
		return nil
	}

	requests := []ctrl.Request{}
	for _, host := range bmhHostList.Items {
		log.Info("found host", "name", host.Name)
		requests = append(requests, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      host.Name,
				Namespace: host.Namespace,
			},
		})
	}
	return requests
}

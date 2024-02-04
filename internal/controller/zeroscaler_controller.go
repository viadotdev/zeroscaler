/*
Copyright 2024 - Jimmy Lipham & Zeroscaler Authors

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	operatorv1 "zeroscaler.com/core/api/v1"
	"zeroscaler.com/core/internal/proxy"
	scaler "zeroscaler.com/core/internal/scaler"
)

// ZeroScalerReconciler reconciles a ZeroScaler object
type ZeroScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	scaler *scaler.Scaler
	proxy  *proxy.Proxy
}

func NewZeroScalerReconciler(client client.Client, scheme *runtime.Scheme) *ZeroScalerReconciler {
	scaler := scaler.NewScaler(client)
	return &ZeroScalerReconciler{
		Client: client,
		Scheme: scheme,
		scaler: scaler,
		proxy:  proxy.NewProxy(context.Background(), scaler),
	}
}

//+kubebuilder:rbac:groups=operator.zeroscaler.com,resources=zeroscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.zeroscaler.com,resources=zeroscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.zeroscaler.com,resources=zeroscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ZeroScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *ZeroScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here
	scaler := &operatorv1.ZeroScaler{}
	if err := r.Get(ctx, req.NamespacedName, scaler); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	l.Info("Found ZeroScaler", "namespace", scaler.Namespace, "name", scaler.Name, "deployment", scaler.Spec.Target.Name, "scaleDownAfter", scaler.Spec.ScaleDownAfter)

	// Check the target deployment
	// Check the last time the deployment was accessed
	// If the time is greater than the scale down time, scale down the deployment
	var deployment appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKey{Name: scaler.Spec.Target.Name, Namespace: req.Namespace}, &deployment); err != nil {
		return ctrl.Result{}, err
	}

	// Check the target service
	var service corev1.Service
	if err := r.Get(ctx, client.ObjectKey{Name: scaler.Spec.Target.Service, Namespace: req.Namespace}, &service); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Found Deployment", "namespace", deployment.Namespace, "name", deployment.Name)

	if service.Spec.Type == corev1.ServiceTypeClusterIP {
		l.Info("Service is of type ClusterIP", "namespace", service.Namespace, "name", service.Name, "port", scaler.Spec.Target.Port)
		target := fmt.Sprintf("http://%s:%d", service.Spec.ClusterIP, scaler.Spec.Target.Port)

		for _, host := range scaler.Spec.Hosts {
			l.Info("Adding route", "host", host, "deployment", deployment.Name, "target", target)
			r.proxy.AddBackend(host, target, &deployment, scaler.Spec.ScaleDownAfter)
		}
	} else if service.Spec.Type == corev1.ServiceTypeNodePort {
		l.Info("Service is of type NodePort", "namespace", service.Namespace, "name", service.Name, "port", scaler.Spec.Target.Port)
	} else if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		l.Info("Service is of type LoadBalancer", "namespace", service.Namespace, "name", service.Name, "port", scaler.Spec.Target.Port)
	}

	// Scale down the deployment
	/*
		deployment.Spec.Replicas = new(int32)
		*deployment.Spec.Replicas = 0
		if err := r.Update(ctx, &deployment); err != nil {
			return ctrl.Result{}, err
		}
	*/

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZeroScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	go r.proxy.Start()
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.ZeroScaler{}).
		Complete(r)
}

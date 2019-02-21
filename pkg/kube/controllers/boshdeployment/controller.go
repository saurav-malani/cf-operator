package boshdeployment

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bdm "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	bdc "code.cloudfoundry.org/cf-operator/pkg/kube/apis/boshdeployment/v1alpha1"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/context"
)

// Add creates a new BOSHDeployment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(log *zap.SugaredLogger, ctrConfig *context.Config, mgr manager.Manager) error {
	r := NewReconciler(log, ctrConfig, mgr, bdm.NewResolver(mgr.GetClient(), bdm.NewInterpolator()), controllerutil.SetControllerReference)

	// Create a new controller
	c, err := controller.New("boshdeployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource BOSHDeployment
	err = c.Watch(&source.Kind{Type: &bdc.BOSHDeployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner BOSHDeployment
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &bdc.BOSHDeployment{},
	})
	if err != nil {
		return err
	}

	return nil
}

package custompod

import (
	"context"
	"strconv"

	bdm "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new BOSHDeployment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(log *zap.SugaredLogger, mgr manager.Manager) error {
	return add(mgr, NewReconciler(log, mgr, controllerutil.SetControllerReference))
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(log *zap.SugaredLogger, mgr manager.Manager, srf setReferenceFunc) reconcile.Reconciler {
	return &ReconcileCustomPod{
		log:          log,
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		setReference: srf,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("custompod-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// watch only our pods
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}, predicate.Funcs{
		CreateFunc: func(ev event.CreateEvent) bool {
			annotations := ev.Meta.GetAnnotations()
			for key := range annotations {
				if key == "custompod" {
					return true
				}
			}
			r.(*ReconcileCustomPod).log.Info("not our pod, ignore create")
			return false
		},
		GenericFunc: func(ev event.GenericEvent) bool {
			return false
		},
		UpdateFunc: func(ev event.UpdateEvent) bool {
			annotations := ev.MetaNew.GetAnnotations()
			for key := range annotations {
				if key == "custompod" {
					return true
				}
			}
			r.(*ReconcileCustomPod).log.Info("not our pod, ignore update")
			return false
		},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCustomPod{}

type setReferenceFunc func(owner, object metav1.Object, scheme *runtime.Scheme) error

type ReconcileCustomPod struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	scheme       *runtime.Scheme
	setReference setReferenceFunc
	log          *zap.SugaredLogger
}

// Reconcile reads that state of the cluster and makes changes based on the state read
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomPod) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.log.Infof("Reconciling Custom Pod %s\n", request.Name)

	instance := &corev1.Pod{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			r.log.Debugf("custom pod controller delete triggered for %s", request.Name)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp != nil {
		r.log.Debugf("custom pod controller delete triggered, remove our finalizer on %s", request.Name)
		instance.ObjectMeta.Finalizers = []string{}
		r.client.Update(context.TODO(), instance)
		return reconcile.Result{}, nil
	}

	r.log.Debugf("custom pod controller add triggered for %s", request.Name)

	// Set instance as the owner and controller of generated objects:
	// if err := r.setReference(instance, ...

	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *bdm.Manifest, namespace string) *corev1.Pod {
	ig := cr.InstanceGroups[0]
	labels := map[string]string{
		"app":  ig.Name,
		"size": strconv.Itoa(ig.Instances),
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ig.Name + "-pod",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

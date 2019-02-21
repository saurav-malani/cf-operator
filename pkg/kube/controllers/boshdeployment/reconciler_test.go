package boshdeployment_test

import (
	"fmt"
	"time"

	bdm "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest/fakes"
	bdc "code.cloudfoundry.org/cf-operator/pkg/kube/apis/boshdeployment/v1alpha1"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers"
	cfd "code.cloudfoundry.org/cf-operator/pkg/kube/controllers/boshdeployment"
	cfakes "code.cloudfoundry.org/cf-operator/pkg/kube/controllers/fakes"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReconcileBoshDeployment", func() {
	var (
		recorder   *record.FakeRecorder
		manager    *cfakes.FakeManager
		reconciler reconcile.Reconciler
		request    reconcile.Request
		resolver   fakes.FakeResolver
		manifest   *bdm.Manifest
		log        *zap.SugaredLogger
		ctrsConfig *context.Config
	)

	BeforeEach(func() {
		controllers.AddToScheme(scheme.Scheme)
		recorder = record.NewFakeRecorder(20)
		manager = &cfakes.FakeManager{}
		manager.GetRecorderReturns(recorder)
		resolver = fakes.FakeResolver{}

		request = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
		manifest = &bdm.Manifest{
			InstanceGroups: []*bdm.InstanceGroup{
				{Name: "fakepod"},
			},
		}
		core, _ := observer.New(zapcore.InfoLevel)
		ctrsConfig = &context.Config{ //Set the context to be TODO
			CtxTimeOut: 10 * time.Second,
			CtxType:    context.NewContext(),
		}
		log = zap.New(core).Sugar()
	})

	JustBeforeEach(func() {
		resolver.ResolveCRDReturns(manifest, nil)
		reconciler = cfd.NewReconciler(log, ctrsConfig, manager, &resolver, controllerutil.SetControllerReference)
	})

	Describe("Reconcile", func() {
		Context("when the manifest can not be resolved", func() {
			var (
				client *cfakes.FakeClient
			)
			BeforeEach(func() {
				client = &cfakes.FakeClient{}
				manager.GetClientReturns(client)
			})

			It("returns an empty Result when the resource was not found", func() {
				client.GetReturns(errors.NewNotFound(schema.GroupResource{}, "not found is requeued"))

				reconciler.Reconcile(request)
				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(reconcile.Result{}).To(Equal(result))
			})

			It("throws an error when the request failed", func() {
				client.GetReturns(errors.NewBadRequest("bad request returns error"))

				_, err := reconciler.Reconcile(request)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("bad request returns error"))

				// check for events
				Expect(<-recorder.Events).To(ContainSubstring("GetCRD Error"))
			})

			It("handles errors when resolving the CR", func() {
				resolver.ResolveCRDReturns(nil, fmt.Errorf("resolver error"))

				_, err := reconciler.Reconcile(request)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("resolver error"))

				// check for events
				Expect(<-recorder.Events).To(ContainSubstring("ResolveCRD Error"))
			})

			It("handles errors when missing instance groups", func() {
				resolver.ResolveCRDReturns(&bdm.Manifest{
					InstanceGroups: []*bdm.InstanceGroup{},
				}, nil)

				_, err := reconciler.Reconcile(request)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("manifest is missing instance groups"))

				// check for events
				Expect(<-recorder.Events).To(ContainSubstring("MissingInstance Error"))
			})
		})

		Context("when the manifest can be resolved", func() {
			var (
				client client.Client
			)
			BeforeEach(func() {
				client = fake.NewFakeClient(
					&bdc.BOSHDeployment{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "foo",
							Namespace: "default",
						},
						Spec: bdc.BOSHDeploymentSpec{},
					},
				)
				manager.GetClientReturns(client)
			})

			Context("With an empty manifest", func() {
				BeforeEach(func() {
					manifest = &bdm.Manifest{}
				})

				It("raises an error if there are no instance groups defined in the manifest", func() {
					_, err := reconciler.Reconcile(request)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("manifest is missing instance groups"))
				})
			})

			It("handles errors when setting the owner reference on the object", func() {
				ctrsConfig := &context.Config{ //Set the context to be TODO
					CtxTimeOut: 10 * time.Second,
					CtxType:    context.NewContext(),
				}
				reconciler = cfd.NewReconciler(log, ctrsConfig, manager, &resolver, func(owner, object metav1.Object, scheme *runtime.Scheme) error {
					return fmt.Errorf("failed to set reference")
				})

				_, err := reconciler.Reconcile(request)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to set reference"))
			})
		})
	})
})

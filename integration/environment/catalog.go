package environment

import (
	fisv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/boshdeploymentcontroller/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Catalog provides several instances for tests
type Catalog struct{}

var terminationGracePeriodSeconds = int64(1)

// DefaultBOSHManifest for tests
func (c *Catalog) DefaultBOSHManifest(name string) corev1.ConfigMap {
	return corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Data: map[string]string{
			"manifest": `instance-groups:
- name: diego
  instances: 3
- name: mysql
`,
		},
	}
}

// DefaultSecret for tests
func (c *Catalog) DefaultSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: v1.ObjectMeta{Name: name},
		StringData: map[string]string{},
	}
}

// AnnotatedSecret for tests
func (c *Catalog) AnnotatedSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"custompod": "yes",
			},
		},
		StringData: map[string]string{},
	}
}

// DefaultPod for tests
func (c *Catalog) DefaultPod(name string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
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

// CustomPod for tests
func (c *Catalog) CustomPod(name string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:       name,
			Finalizers: []string{"fissile.app/finalizer-test"},
			Annotations: map[string]string{
				"custompod": "yes",
			},
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
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

// DefaultFissileCR fissile deployment CR
func (c *Catalog) DefaultFissileCR(name, manifestRef string) fisv1.BOSHDeployment {
	return fisv1.BOSHDeployment{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Spec: fisv1.BOSHDeploymentSpec{
			ManifestRef: manifestRef,
		},
	}
}

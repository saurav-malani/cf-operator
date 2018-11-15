package manifest_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"code.cloudfoundry.org/fissile/model"
	"code.cloudfoundry.org/fissile/model/resolver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DummyResolver struct {
	releases model.Releases
}

func NewDummyResolver() *DummyResolver {
	return &DummyResolver{}
}

func (d *DummyResolver) Load(options model.ReleaseOptions, releaseRefs []*model.ReleaseRef) (model.Releases, error) {
	return d.releases, nil
}

func (d *DummyResolver) CanValidate() bool {
	return true
}

func (d *DummyResolver) MapReleases(r model.Releases) error {
	return nil
}

func (d *DummyResolver) FindRelease(name string) (*model.Release, bool) {
	// this will make cf-operator tests pass
	return &model.Release{
		Name: name,
		Jobs: model.Jobs{
			&model.Job{Name: name},
			&model.Job{Name: "new_hostname"},
		},
	}, true
}

func LoadRoleManifest(manifestFilePath string, options model.LoadRoleManifestOptions) (*model.RoleManifest, error) {
	roleManifest := model.NewRoleManifest()
	err := roleManifest.LoadManifestFromFile(manifestFilePath)
	if err != nil {
		return nil, err
	}

	return resolver.NewResolver(roleManifest, NewDummyResolver(), options).Resolve()
}

var _ = Describe("LoadManifest", func() {
	var (
		manifest *model.RoleManifest
		tmpFile  *os.File
		err      error
	)

	BeforeEach(func() {
	})

	FIt("loads empty manifest", func() {
		tmpFile, err = ioutil.TempFile("", "bosh-test-manifest")
		Expect(err).ToNot(HaveOccurred())
		manifest, err = LoadRoleManifest(tmpFile.Name(), model.LoadRoleManifestOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(manifest).ToNot(Equal(nil))
		Expect(len(manifest.InstanceGroups)).To(Equal(0))
		defer os.Remove(tmpFile.Name())
	})

	FIt("loads valid manifest", func() {
		workDir, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		manifestPath := filepath.Join(workDir, "../../../../fissile/test-assets/role-manifests/model/tor-good.yml")

		manifest, err = LoadRoleManifest(manifestPath, model.LoadRoleManifestOptions{})
		Expect(err).ToNot(HaveOccurred())
		d, _ := yaml.Marshal(manifest)
		fmt.Printf("%s\n", string(d))

		Expect(manifest).ToNot(Equal(nil))
		Expect(len(manifest.InstanceGroups)).To(Equal(2))
		defer os.Remove(tmpFile.Name())
	})
})

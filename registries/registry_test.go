//
// Copyright (c) 2018 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package registries

import (
	"fmt"
	"testing"

	"github.com/automationbroker/bundle-lib/bundle"
	"github.com/automationbroker/bundle-lib/registries/adapters"
	"github.com/stretchr/testify/assert"
)

var SpecTags = []string{"latest", "old-release"}

const SpecID = "ab094014-b740-495e-b178-946d5aa97ebf"
const SpecBadVersion = "2.0.0"
const SpecVersion = "1.0.0"
const SpecRuntime = 1
const SpecBadRuntime = 0
const SpecName = "etherpad-bundle"
const SpecImage = "fusor/etherpad-bundle"
const SpecBindable = false
const SpecAsync = "optional"
const SpecDescription = "A note taking webapp"
const SpecRegistryName = "test"

const PlanName = "dev"
const PlanDescription = "Basic development plan"

var PlanMetadata = map[string]interface{}{
	"displayName":     "Development",
	"longDescription": PlanDescription,
	"cost":            "$0.00",
}

const PlanFree = true
const PlanBindable = true

var expectedPlanParameters = []bundle.ParameterDescriptor{
	{
		Name:    "postgresql_database",
		Default: "admin",
		Type:    "string",
		Title:   "PostgreSQL Database Name",
	},
	{
		Name:        "postgresql_password",
		Default:     "admin",
		Type:        "string",
		Description: "A random alphanumeric string if left blank",
		Title:       "PostgreSQL Password",
	},
	{
		Name:                "postgresql_user",
		Default:             "admin",
		Title:               "PostgreSQL User",
		Type:                "string",
		DeprecatedMaxlength: 63,
	},
	{
		Name:    "postgresql_version",
		Default: 9.5,
		Enum:    []string{"9.5", "9.4"},
		Type:    "enum",
		Title:   "PostgreSQL Version",
	},
	{
		Name:        "postgresql_email",
		Pattern:     "\u201c^\\\\S+@\\\\S+$\u201d",
		Type:        "string",
		Description: "email address",
		Title:       "email",
	},
}

var p = bundle.Plan{
	Name:        PlanName,
	Description: PlanDescription,
	Metadata:    PlanMetadata,
	Free:        PlanFree,
	Bindable:    PlanBindable,
	Parameters:  expectedPlanParameters,
}

var s = bundle.Spec{
	Version:     SpecVersion,
	Runtime:     SpecRuntime,
	ID:          SpecID,
	Description: SpecDescription,
	FQName:      SpecName,
	Image:       SpecImage,
	Tags:        SpecTags,
	Bindable:    SpecBindable,
	Async:       SpecAsync,
	Plans:       []bundle.Plan{p},
}

var noPlansSpec = bundle.Spec{
	Version:     SpecVersion,
	Runtime:     SpecRuntime,
	ID:          SpecID,
	Description: SpecDescription,
	FQName:      SpecName,
	Image:       SpecImage,
	Tags:        SpecTags,
	Bindable:    SpecBindable,
	Async:       SpecAsync,
}

var noVersionSpec = bundle.Spec{
	Runtime:     SpecRuntime,
	ID:          SpecID,
	Description: SpecDescription,
	FQName:      SpecName,
	Image:       SpecImage,
	Tags:        SpecTags,
	Bindable:    SpecBindable,
	Async:       SpecAsync,
	Plans:       []bundle.Plan{p},
}

var badVersionSpec = bundle.Spec{
	Version:     SpecBadVersion,
	Runtime:     SpecRuntime,
	ID:          SpecID,
	Description: SpecDescription,
	FQName:      SpecName,
	Image:       SpecImage,
	Tags:        SpecTags,
	Bindable:    SpecBindable,
	Async:       SpecAsync,
	Plans:       []bundle.Plan{p},
}

var badRuntimeSpec = bundle.Spec{
	Version:     SpecVersion,
	Runtime:     SpecBadRuntime,
	ID:          SpecID,
	Description: SpecDescription,
	FQName:      SpecName,
	Image:       SpecImage,
	Tags:        SpecTags,
	Bindable:    SpecBindable,
	Async:       SpecAsync,
	Plans:       []bundle.Plan{p},
}

type TestingAdapter struct {
	Name   string
	Images []string
	Specs  []*bundle.Spec
	Called map[string]bool
}

func (t TestingAdapter) GetImageNames() ([]string, error) {
	t.Called["GetImageNames"] = true
	return t.Images, nil
}

func (t TestingAdapter) FetchSpecs(images []string) ([]*bundle.Spec, error) {
	t.Called["FetchSpecs"] = true
	return t.Specs, nil
}

func (t TestingAdapter) RegistryName() string {
	t.Called["RegistryName"] = true
	return t.Name
}

var a *TestingAdapter
var r Registry

func setUp() Registry {
	a = &TestingAdapter{
		Name:   "testing",
		Images: []string{"image1-bundle", "image2"},
		Specs:  []*bundle.Spec{&s},
		Called: map[string]bool{},
	}
	filter := Filter{}
	c := Config{}
	r = Registry{config: c,
		adapter: a,
		filter:  filter}
	return r
}

func setUpNoPlans() Registry {
	a = &TestingAdapter{
		Name:   "testing",
		Images: []string{"image1-bundle", "image2"},
		Specs:  []*bundle.Spec{&noPlansSpec},
		Called: map[string]bool{},
	}
	filter := Filter{}
	c := Config{}
	r = Registry{config: c,
		adapter: a,
		filter:  filter}
	return r
}

func setUpNoVersion() Registry {
	a = &TestingAdapter{
		Name:   "testing",
		Images: []string{"image1-bundle", "image2"},
		Specs:  []*bundle.Spec{&noVersionSpec},
		Called: map[string]bool{},
	}
	filter := Filter{}
	c := Config{}
	r = Registry{config: c,
		adapter: a,
		filter:  filter}
	return r
}

func setUpBadVersion() Registry {
	a = &TestingAdapter{
		Name:   "testing",
		Images: []string{"image1-bundle", "image2"},
		Specs:  []*bundle.Spec{&badVersionSpec},
		Called: map[string]bool{},
	}
	filter := Filter{}
	c := Config{}
	r = Registry{config: c,
		adapter: a,
		filter:  filter}
	return r
}

func setUpBadRuntime() Registry {
	a = &TestingAdapter{
		Name:   "testing",
		Images: []string{"image1-bundle", "image2"},
		Specs:  []*bundle.Spec{&badRuntimeSpec},
		Called: map[string]bool{},
	}
	filter := Filter{}
	c := Config{}
	r = Registry{config: c,
		adapter: a,
		filter:  filter}
	return r
}

func TestRegistryLoadSpecsNoError(t *testing.T) {
	r := setUp()
	specs, numImages, err := r.LoadSpecs()
	if err != nil {
		assert.True(t, false)
	}
	assert.True(t, a.Called["GetImageNames"])
	assert.True(t, a.Called["FetchSpecs"])
	assert.Equal(t, numImages, 2)
	assert.Equal(t, len(specs), 1)
	assert.Equal(t, specs[0], &s)
}

func TestRegistryLoadSpecsNoPlans(t *testing.T) {
	r := setUpNoPlans()
	specs, _, err := r.LoadSpecs()
	if err != nil {
		assert.True(t, false)
	}
	assert.True(t, a.Called["GetImageNames"])
	assert.True(t, a.Called["FetchSpecs"])
	assert.Equal(t, len(specs), 0)
}

func TestRegistryLoadSpecsNoVersion(t *testing.T) {
	r := setUpNoVersion()
	specs, _, err := r.LoadSpecs()
	if err != nil {
		assert.True(t, false)
	}
	assert.True(t, a.Called["GetImageNames"])
	assert.True(t, a.Called["FetchSpecs"])
	assert.Equal(t, len(specs), 0)
}

func TestRegistryLoadSpecsBadVersion(t *testing.T) {
	r := setUpBadVersion()
	specs, _, err := r.LoadSpecs()
	if err != nil {
		assert.True(t, false)
	}
	assert.True(t, a.Called["GetImageNames"])
	assert.True(t, a.Called["FetchSpecs"])
	assert.Equal(t, len(specs), 0)
}

func TestRegistryLoadSpecsBadRuntime(t *testing.T) {
	r := setUpBadRuntime()
	specs, _, err := r.LoadSpecs()
	if err != nil {
		assert.True(t, false)
	}
	assert.True(t, a.Called["GetImageNames"])
	assert.True(t, a.Called["FetchSpecs"])
	assert.Equal(t, len(specs), 0)
}

func TestFail(t *testing.T) {
	inputerr := fmt.Errorf("sample test err")

	testCases := []struct {
		name     string
		r        Registry
		expected bool
	}{
		{
			name: "fail should return true",
			r: Registry{
				config: Config{
					Fail: true,
				},
			},
			expected: true,
		},
		{
			name: "fail should return false",
			r: Registry{
				config: Config{
					Fail: false,
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.r.Fail(inputerr))
		})
	}
}

func TestNewRegistryRHCC(t *testing.T) {
	c := Config{
		Type: "rhcc",
		Name: "rhcc",
	}
	reg, err := NewRegistry(c, "")
	if err != nil {
		assert.True(t, false)
	}
	_, ok := reg.adapter.(*adapters.RHCCAdapter)
	assert.True(t, ok)
}

func TestNewRegistryDockerHub(t *testing.T) {
	c := Config{
		Type: "dockerhub",
		Name: "dh",
		URL:  "https://registry.hub.docker.com",
		User: "shurley",
		Org:  "shurley",
	}
	reg, err := NewRegistry(c, "")
	if err != nil {
		assert.True(t, false)
	}
	_, ok := reg.adapter.(*adapters.DockerHubAdapter)
	assert.True(t, ok)
}

func TestNewRegistryMock(t *testing.T) {
	c := Config{
		Type: "mock",
		Name: "mock",
	}

	reg, err := NewRegistry(c, "")
	if err != nil {
		assert.True(t, false)
	}
	_, ok := reg.adapter.(*adapters.MockAdapter)
	assert.True(t, ok)
}

func TestUnknownType(t *testing.T) {
	c := Config{
		Type: "makes_no_sense",
		Name: "dh",
	}
	_, err := NewRegistry(c, "")
	if err == nil {
		t.Fatal("Error: error was nil")
	}
}

func TestRegistryName(t *testing.T) {
	testCases := []struct {
		name     string
		r        Registry
		expected string
	}{
		{
			name: "registry name",
			r: Registry{
				config: Config{
					Name: "registryname",
				},
			},
			expected: "registryname",
		},
		{
			name: "empty name",
			r: Registry{
				config: Config{
					Name: "",
				},
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.r.RegistryName())
		})
	}
}

func TestValidate(t *testing.T) {

	testCases := []struct {
		name     string
		c        Config
		expected bool
	}{
		{
			name:     "empty name",
			c:        Config{Name: ""},
			expected: false,
		},
		{
			name: "valid name, empty authtype and authname",
			c: Config{
				Name:     "daname",
				AuthName: "",
				AuthType: "",
			},
			expected: true,
		},
		{
			name: "valid name, empty authtype, non-empty authname",
			c: Config{
				Name:     "daname",
				AuthName: "shouldfail",
				AuthType: "",
			},
			expected: false,
		},
		{
			name: "valid name, file, empty authname",
			c: Config{
				Name:     "daname",
				AuthName: "",
				AuthType: "file",
			},
			expected: false,
		},
		{
			name: "valid name, file, non-empty authname",
			c: Config{
				Name:     "daname",
				AuthName: "non-empty",
				AuthType: "file",
			},
			expected: true,
		},
		{
			name: "valid name, secret, empty authname",
			c: Config{
				Name:     "daname",
				AuthName: "",
				AuthType: "secret",
			},
			expected: false,
		},
		{
			name: "valid name, secret, non-empty authname",
			c: Config{
				Name:     "daname",
				AuthName: "non-empty",
				AuthType: "secret",
			},
			expected: true,
		},
		{
			name: "valid name, config, without user",
			c: Config{
				Name:     "daname",
				User:     "",
				AuthType: "config",
			},
			expected: false,
		},
		{
			name: "valid name, config, without pass",
			c: Config{
				Name:     "daname",
				User:     "user",
				Pass:     "",
				AuthType: "config",
			},
			expected: false,
		},
		{
			name: "valid name, config, user, pass",
			c: Config{
				Name:     "daname",
				User:     "user",
				Pass:     "$3cr3+",
				AuthType: "config",
			},
			expected: true,
		},
		{
			name: "valid name, unknown",
			c: Config{
				Name:     "daname",
				AuthType: "unknown",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.c.Validate())
		})
	}
}

type fakeAdapter struct{}

func (f fakeAdapter) GetImageNames() ([]string, error) {
	return []string{}, nil
}

func (f fakeAdapter) FetchSpecs(names []string) ([]*bundle.Spec, error) {
	return []*bundle.Spec{}, nil
}

func (f fakeAdapter) RegistryName() string {
	return ""
}

func TestAdapterWithConfiguration(t *testing.T) {
	c := Config{
		Name: "nsa",
		Type: "custom",
	}

	f := fakeAdapter{}

	reg, err := NewCustomRegistry(c, f, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, reg.adapter, f, "registry uses wrong adapter")
	assert.Equal(t, reg.config, c, "registrying using wrong config")
}

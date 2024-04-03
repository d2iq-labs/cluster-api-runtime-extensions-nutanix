package instancetype

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInstanceTypePatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InstanceType patches for ControlPlane and Workers suite")
}

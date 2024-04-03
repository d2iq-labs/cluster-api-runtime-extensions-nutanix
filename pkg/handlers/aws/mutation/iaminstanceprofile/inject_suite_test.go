package iaminstanceprofile

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIAMInstnaceProfilePatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IAMInstanceProfile patches for ControlPlane and Workers suite")
}

package region

import (
	"testing"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers/mutation"
	regiontests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/aws/mutation/region/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRegionPatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS region patch suite")
}

var _ = Describe("AWS region patches", func() {
	// only add aws region patch
	regionPatchGenerator := func() mutation.GeneratePatches {
		return mutation.NewMetaGeneratePatchesHandler("", NewPatch()).(mutation.GeneratePatches)
	}
	regiontests.TestGeneratePatches(
		GinkgoT(),
		regionPatchGenerator,
		clusterconfig.MetaVariableName,
		v1alpha1.AWSVariableName,
		VariableName,
	)
})

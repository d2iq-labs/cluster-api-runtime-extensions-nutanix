// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers/mutation"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/testutils/capitest"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/testutils/capitest/request"
	"github.com/onsi/gomega"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"
)

func TestGeneratePatches(
	t *testing.T,
	generatorFunc func() mutation.GeneratePatches,
	variableName string,
	variablePath ...string,
) {
	t.Helper()

	capitest.ValidateGeneratePatches(
		t,
		generatorFunc,
		capitest.PatchTestDef{
			Name: "unset variable",
		},
		capitest.PatchTestDef{
			Name: "ControlPlaneEndpoint set to valid host",
			Vars: []runtimehooksv1.Variable{
				capitest.VariableWithValue(
					variableName,
					v1alpha1.NutanixControlPlaneEndpointSpec{
						Host: "10.20.100.10",
						Port: 6443,
					},
					variablePath...,
				),
			},
			RequestItem: request.NewNutanixClusterTemplateRequestItem("1"),
			ExpectedPatchMatchers: []capitest.JSONPatchMatcher{
				{
					Operation:    "replace",
					Path:         "/spec/template/spec/controlPlaneEndpoint/host",
					ValueMatcher: gomega.Equal("10.20.100.10"),
				},
			},
		},
		capitest.PatchTestDef{
			Name: "ControlPlaneEndpoint set to valid host",
			Vars: []runtimehooksv1.Variable{
				capitest.VariableWithValue(
					variableName,
					v1alpha1.NutanixControlPlaneEndpointSpec{
						Host: "10.20.100.10",
						Port: 6443,
					},
					variablePath...,
				),
			},
			RequestItem: request.NewNutanixClusterTemplateRequestItem("2"),
			ExpectedPatchMatchers: []capitest.JSONPatchMatcher{
				{
					Operation:    "replace",
					Path:         "/spec/template/spec/controlPlaneEndpoint/port",
					ValueMatcher: gomega.BeEquivalentTo(6443),
				},
			},
		},
	)
}

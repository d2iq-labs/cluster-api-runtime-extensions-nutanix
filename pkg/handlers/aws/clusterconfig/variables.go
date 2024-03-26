// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package clusterconfig

import (
	"context"

	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	commonhandlers "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers/mutation"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
)

var (
	_ commonhandlers.Named       = &awsClusterConfigVariableHandler{}
	_ mutation.DiscoverVariables = &awsClusterConfigVariableHandler{}
)

const (
	// HandlerNameVariable is the name of the variable handler.
	HandlerNameVariable = "AWSClusterConfigVars"
)

func NewVariable() *awsClusterConfigVariableHandler {
	return &awsClusterConfigVariableHandler{}
}

type awsClusterConfigVariableHandler struct{}

func (h *awsClusterConfigVariableHandler) Name() string {
	return HandlerNameVariable
}

func (h *awsClusterConfigVariableHandler) DiscoverVariables(
	ctx context.Context,
	_ *runtimehooksv1.DiscoverVariablesRequest,
	resp *runtimehooksv1.DiscoverVariablesResponse,
) {
	resp.Variables = append(resp.Variables, clusterv1.ClusterClassVariable{
		Name:     clusterconfig.MetaVariableName,
		Required: true,
		Schema: v1alpha1.AWSClusterConfigSpec{
			AWS: v1alpha1.AWSSpec{},
			ControlPlane: &v1alpha1.NodeConfigSpec{
				AWS: &v1alpha1.AWSNodeSpec{
					InstanceType: ptr.To(v1alpha1.InstanceType("m5.large")),
				},
			},
		}.VariableSchema(),
	})
	resp.SetStatus(runtimehooksv1.ResponseStatusSuccess)
}

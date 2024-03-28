// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package csi

import (
	"testing"

	"k8s.io/utils/ptr"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/testutils/capitest"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
)

func TestVariableValidation(t *testing.T) {
	capitest.ValidateDiscoverVariables(
		t,
		clusterconfig.MetaVariableName,
		ptr.To(v1alpha1.GenericClusterConfig{}.VariableSchema()),
		false,
		clusterconfig.NewVariable,
		capitest.VariableTestDef{
			Name: "set with empty CSI",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with with empty CSI providers slice",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with invalid provider",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "csi-provider",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
						}},
					},
				},
			},
			ExpectError: true,
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single valid provider with HelmAddon strategy",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyHelmAddon,
						}},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single valid provider with CRS strategy",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
						}},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single valid provider with empty storage class config",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:               "aws-ebs",
							Strategy:           v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{},
						}},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single invalid provider with missing name",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:               "aws-ebs",
							Strategy:           v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{{}},
						}},
					},
				},
			},
			ExpectError: true,
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single valid provider using defaults",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{{
								Name: "default",
							}},
						}},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single valid provider with single empty specified storage class config",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{{
								Name:       "default",
								Parameters: map[string]string{},
							}},
						}},
					},
				},
			},
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single invalid provider with invalid reclaim policy",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{{
								Name:          "default",
								ReclaimPolicy: "invalid",
							}},
						}},
					},
				},
			},
			ExpectError: true,
		},
		capitest.VariableTestDef{
			Name: "set with single CSIProvider with single invalid provider with invalid reclaim volume binding mode",
			Vals: v1alpha1.GenericClusterConfig{
				Addons: &v1alpha1.Addons{
					CSI: &v1alpha1.CSI{
						Providers: []v1alpha1.CSIProvider{{
							Name:     "aws-ebs",
							Strategy: v1alpha1.AddonStrategyClusterResourceSet,
							StorageClassConfig: []v1alpha1.StorageClassConfig{{
								Name:              "default",
								VolumeBindingMode: "invalid",
							}},
						}},
					},
				},
			},
			ExpectError: true,
		},
	)
}

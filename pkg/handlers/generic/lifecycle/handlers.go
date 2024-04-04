// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	"github.com/spf13/pflag"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/ccm"
	awsccm "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/ccm/aws"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/clusterautoscaler"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/cni/calico"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/cni/cilium"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/csi"
	awsebs "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/csi/aws-ebs"
	nutanixcsi "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/csi/nutanix"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/nfd"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/servicelbgc"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/options"
)

type Handlers struct {
	calicoCNIConfig         *calico.CNIConfig
	ciliumCNIConfig         *cilium.CNIConfig
	nfdConfig               *nfd.Config
	clusterAutoscalerConfig *clusterautoscaler.Config
	ebsCSIConfig            *awsebs.Config
	nutanixCSIConfig        *nutanixcsi.Config
	awsCCMConfig            *awsccm.Config
}

func New(globalOptions *options.GlobalOptions) *Handlers {
	return &Handlers{
		calicoCNIConfig:         &calico.CNIConfig{GlobalOptions: globalOptions},
		ciliumCNIConfig:         &cilium.CNIConfig{GlobalOptions: globalOptions},
		nfdConfig:               &nfd.Config{GlobalOptions: globalOptions},
		clusterAutoscalerConfig: &clusterautoscaler.Config{GlobalOptions: globalOptions},
		ebsCSIConfig:            &awsebs.Config{GlobalOptions: globalOptions},
		awsCCMConfig:            &awsccm.Config{GlobalOptions: globalOptions},
		nutanixCSIConfig:        &nutanixcsi.Config{GlobalOptions: globalOptions},
	}
}

func (h *Handlers) AllHandlers(mgr manager.Manager) []handlers.Named {
	csiHandlers := map[string]csi.Provider{
		v1alpha1.CSIProviderAWSEBS:  awsebs.New(mgr.GetClient(), h.ebsCSIConfig),
		v1alpha1.CSIProviderNutanix: nutanixcsi.New(mgr.GetClient(), h.nutanixCSIConfig),
	}
	ccmHandlers := map[string]ccm.Provider{
		v1alpha1.CCMProviderAWS: awsccm.New(mgr.GetClient(), h.awsCCMConfig),
	}

	return []handlers.Named{
		calico.New(mgr.GetClient(), h.calicoCNIConfig),
		cilium.New(mgr.GetClient(), h.ciliumCNIConfig),
		nfd.New(mgr.GetClient(), h.nfdConfig),
		clusterautoscaler.New(mgr.GetClient(), h.clusterAutoscalerConfig),
		servicelbgc.New(mgr.GetClient()),
		csi.New(mgr.GetClient(), csiHandlers),
		ccm.New(mgr.GetClient(), ccmHandlers),
	}
}

func (h *Handlers) AddFlags(flagSet *pflag.FlagSet) {
	h.nfdConfig.AddFlags("nfd", flagSet)
	h.clusterAutoscalerConfig.AddFlags("cluster-autoscaler", flagSet)
	h.calicoCNIConfig.AddFlags("cni.calico", flagSet)
	h.ciliumCNIConfig.AddFlags("cni.cilium", flagSet)
	h.ebsCSIConfig.AddFlags("awsebs", pflag.CommandLine)
	h.awsCCMConfig.AddFlags("awsccm", pflag.CommandLine)
	h.nutanixCSIConfig.AddFlags("nutanixcsi", flagSet)
}

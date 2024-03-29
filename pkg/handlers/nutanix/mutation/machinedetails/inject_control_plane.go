// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package machinedetails

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	capxv1 "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/external/github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/patches"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/patches/selectors"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/variables"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
)

const (
	// VariableName is the external patch variable name.
	VariableName = "machineDetails"
)

type nutanixMachineDetailsControlPlanePatchHandler struct {
	variableName      string
	variableFieldPath []string
}

func NewControlPlanePatch() *nutanixMachineDetailsControlPlanePatchHandler {
	return newNutanixMachineDetailsControlPlanePatchHandler(
		clusterconfig.MetaVariableName,
		clusterconfig.MetaControlPlaneConfigName,
		v1alpha1.NutanixVariableName,
		VariableName,
	)
}

func newNutanixMachineDetailsControlPlanePatchHandler(
	variableName string,
	variableFieldPath ...string,
) *nutanixMachineDetailsControlPlanePatchHandler {
	return &nutanixMachineDetailsControlPlanePatchHandler{
		variableName:      variableName,
		variableFieldPath: variableFieldPath,
	}
}

func (h *nutanixMachineDetailsControlPlanePatchHandler) Mutate(
	ctx context.Context,
	obj *unstructured.Unstructured,
	vars map[string]apiextensionsv1.JSON,
	holderRef runtimehooksv1.HolderReference,
	_ client.ObjectKey,
) error {
	log := ctrl.LoggerFrom(ctx).WithValues(
		"holderRef", holderRef,
	)

	nutanixNode, found, err := variables.Get[v1alpha1.NutanixNodeSpec](
		vars,
		h.variableName,
		h.variableFieldPath...,
	)
	if err != nil {
		return err
	}
	if !found {
		log.V(5).Info("Nutanix machine details variable for control-plane not defined")
		return nil
	}

	log = log.WithValues(
		"variableName",
		h.variableName,
		"variableFieldPath",
		h.variableFieldPath,
		"variableValue",
		nutanixNode,
	)

	return patches.MutateIfApplicable(
		obj,
		vars,
		&holderRef,
		selectors.InfrastructureControlPlaneMachines(
			"v1beta1",
			"NutanixMachineTemplate",
		),
		log,
		func(obj *capxv1.NutanixMachineTemplate) error {
			log.WithValues(
				"patchedObjectKind", obj.GetObjectKind().GroupVersionKind().String(),
				"patchedObjectName", client.ObjectKeyFromObject(obj),
			).Info("setting Nutanix machine details in control plane NutanixMachineTemplate spec")

			obj.Spec.Template.Spec.BootType = capxv1.NutanixBootType(
				nutanixNode.BootType,
			)
			obj.Spec.Template.Spec.Cluster = capxv1.NutanixResourceIdentifier{
				Type: nutanixNode.Cluster.Type,
			}
			if nutanixNode.Cluster.Type == capxv1.NutanixIdentifierName {
				obj.Spec.Template.Spec.Cluster.Name = nutanixNode.Cluster.Name
			} else {
				obj.Spec.Template.Spec.Cluster.UUID = nutanixNode.Cluster.UUID
			}

			obj.Spec.Template.Spec.Image = capxv1.NutanixResourceIdentifier{
				Type: nutanixNode.Image.Type,
			}
			if nutanixNode.Image.Type == capxv1.NutanixIdentifierName {
				obj.Spec.Template.Spec.Image.Name = nutanixNode.Image.Name
			} else {
				obj.Spec.Template.Spec.Image.UUID = nutanixNode.Image.UUID
			}

			obj.Spec.Template.Spec.VCPUSockets = nutanixNode.VCPUSockets
			obj.Spec.Template.Spec.VCPUsPerSocket = nutanixNode.VCPUsPerSocket
			obj.Spec.Template.Spec.MemorySize = resource.MustParse(
				nutanixNode.MemorySize,
			)
			obj.Spec.Template.Spec.SystemDiskSize = resource.MustParse(
				nutanixNode.SystemDiskSize,
			)

			subnets := make(
				[]capxv1.NutanixResourceIdentifier,
				len(nutanixNode.Subnets),
			)
			for _, subnetCRE := range nutanixNode.Subnets {
				subnet := capxv1.NutanixResourceIdentifier{
					Type: subnetCRE.Type,
				}
				if subnetCRE.Type == capxv1.NutanixIdentifierName {
					subnet.Name = subnetCRE.Name
				} else {
					subnet.UUID = subnetCRE.UUID
				}
				subnets = append(subnets, subnet)
			}

			obj.Spec.Template.Spec.Subnets = subnets
			// TODO:deepakm-ntnx uncomment this once we are ready
			// obj.Spec.Template.Spec.Project = nutanixNode.Project
			// obj.Spec.Template.Spec.AdditionalCategories = nutanixNode.AdditionalCategories
			// obj.Spec.Template.Spec.GPUs = nutanixNode.GPUs
			return nil
		},
	)
}

// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package machinedetails

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	capxv1 "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/external/github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/patches"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/variables"
)

const (
	// VariableName is the external patch variable name.
	VariableName = "machineDetails"
)

type nutanixMachineDetailsPatchHandler struct {
	metaVariableName  string
	variableFieldPath []string
	patchSelector     clusterv1.PatchSelector
}

func newNutanixMachineDetailsPatchHandler(
	metaVariableName string,
	variableFieldPath []string,
	patchSelector clusterv1.PatchSelector,
) *nutanixMachineDetailsPatchHandler {
	return &nutanixMachineDetailsPatchHandler{
		metaVariableName:  metaVariableName,
		variableFieldPath: variableFieldPath,
		patchSelector:     patchSelector,
	}
}

func (h *nutanixMachineDetailsPatchHandler) Mutate(
	ctx context.Context,
	obj *unstructured.Unstructured,
	vars map[string]apiextensionsv1.JSON,
	holderRef runtimehooksv1.HolderReference,
	_ client.ObjectKey,
) error {
	log := ctrl.LoggerFrom(ctx).WithValues(
		"holderRef", holderRef,
	)

	nutanixMachineDetailsVar, found, err := variables.Get[v1alpha1.NutanixMachineDetails](
		vars,
		h.metaVariableName,
		h.variableFieldPath...,
	)
	if err != nil {
		return err
	}
	if !found {
		log.V(5).Info("Nutanix machine details variable for workers not defined")
		return nil
	}

	log = log.WithValues(
		"variableName",
		h.metaVariableName,
		"variableFieldPath",
		h.variableFieldPath,
		"variableValue",
		nutanixMachineDetailsVar,
	)

	return patches.MutateIfApplicable(
		obj,
		vars,
		&holderRef,
		h.patchSelector,
		log,
		func(obj *capxv1.NutanixMachineTemplate) error {
			log.WithValues(
				"patchedObjectKind", obj.GetObjectKind().GroupVersionKind().String(),
				"patchedObjectName", client.ObjectKeyFromObject(obj),
			).Info("setting Nutanix machine details in worker NutanixMachineTemplate spec")

			obj.Spec.Template.Spec.BootType = capxv1.NutanixBootType(
				nutanixMachineDetailsVar.BootType,
			)
			obj.Spec.Template.Spec.Cluster = capxv1.NutanixResourceIdentifier{
				Type: capxv1.NutanixIdentifierType(nutanixMachineDetailsVar.Cluster.Type),
			}
			if nutanixMachineDetailsVar.Cluster.Type == v1alpha1.NutanixIdentifierName {
				obj.Spec.Template.Spec.Cluster.Name = nutanixMachineDetailsVar.Cluster.Name
			} else {
				obj.Spec.Template.Spec.Cluster.UUID = nutanixMachineDetailsVar.Cluster.UUID
			}

			obj.Spec.Template.Spec.Image = capxv1.NutanixResourceIdentifier{
				Type: capxv1.NutanixIdentifierType(nutanixMachineDetailsVar.Image.Type),
			}
			if nutanixMachineDetailsVar.Image.Type == v1alpha1.NutanixIdentifierName {
				obj.Spec.Template.Spec.Image.Name = nutanixMachineDetailsVar.Image.Name
			} else {
				obj.Spec.Template.Spec.Image.UUID = nutanixMachineDetailsVar.Image.UUID
			}

			obj.Spec.Template.Spec.VCPUSockets = nutanixMachineDetailsVar.VCPUSockets
			obj.Spec.Template.Spec.VCPUsPerSocket = nutanixMachineDetailsVar.VCPUsPerSocket
			obj.Spec.Template.Spec.MemorySize = resource.MustParse(
				nutanixMachineDetailsVar.MemorySize,
			)
			obj.Spec.Template.Spec.SystemDiskSize = resource.MustParse(
				nutanixMachineDetailsVar.SystemDiskSize,
			)
			obj.Spec.Template.Spec.Subnets = make([]capxv1.NutanixResourceIdentifier, 0)
			for _, subnetIdentifier := range nutanixMachineDetailsVar.Subnets {
				if subnetIdentifier.Type == "" {
					continue
				}
				subnet := capxv1.NutanixResourceIdentifier{}
				if subnetIdentifier.Type == v1alpha1.NutanixIdentifierName {
					subnet.Type = capxv1.NutanixIdentifierName
					if subnetIdentifier.Name == nil || *subnetIdentifier.Name == "" {
						continue
					}
					subnet.Name = subnetIdentifier.Name
				} else {
					subnet.Type = capxv1.NutanixIdentifierUUID
					if subnetIdentifier.UUID == nil || *subnetIdentifier.UUID == "" {
						continue
					}
					subnet.UUID = subnetIdentifier.UUID
				}
				obj.Spec.Template.Spec.Subnets = append(obj.Spec.Template.Spec.Subnets, subnet)
			}
			// TODO:deepakm-ntnx assign user provided values
			obj.Spec.Template.Spec.Project = nil
			obj.Spec.Template.Spec.AdditionalCategories = []capxv1.NutanixCategoryIdentifier{}
			obj.Spec.Template.Spec.GPUs = []capxv1.NutanixGPU{}

			return nil
		},
	)
}

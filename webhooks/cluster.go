// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package webhooks

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/variables"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/workerconfig"
)

// +kubebuilder:webhook:verbs=create;update,path=/mutate-cluster-x-k8s-io-v1beta1-cluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=cluster.x-k8s.io,resources=clusters,versions=v1beta1,name=default.cluster.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// Cluster implements a defaulting webhook for Cluster.
type Cluster struct{}

var _ webhook.CustomDefaulter = &Cluster{}

// SetupWebhookWithManager sets up Cluster webhooks.
func (webhook *Cluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&clusterv1.Cluster{}).
		WithDefaulter(webhook).
		Complete()
}

// Default satisfies the defaulting webhook interface.
func (webhook *Cluster) Default(_ context.Context, obj runtime.Object) error {
	// We gather all defaulting errors and return them together.
	var allErrs field.ErrorList

	cluster, ok := obj.(*clusterv1.Cluster)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a Cluster but got a %T", obj))
	}

	if cluster.Spec.Topology == nil ||
		len(cluster.Spec.Topology.Variables) == 0 {
		return nil
	}

	// Set defaults for 'clusterConfig' variable from spec.topology.variables
	clusterConfigVariable, clusterConfigVariableIndex := variables.GetClusterVariableByName(
		clusterconfig.MetaVariableName,
		cluster.Spec.Topology.Variables,
	)
	if clusterConfigVariable != nil {
		clusterConfigSpec := &v1alpha1.ClusterConfigSpec{}
		err := clusterConfigSpec.FromClusterVariable(clusterConfigVariable)
		if err != nil {
			return fmt.Errorf("failed to unmarshal ClusterConfigSpec from ClusterVariable: %w", err)
		}
		errs := defaultClusterConfig(clusterConfigSpec)
		if len(errs) > 1 {
			allErrs = append(allErrs, errs...)
		}
		clusterConfigVariable, err = clusterConfigSpec.ToClusterVariable(clusterconfig.MetaVariableName)
		if err != nil {
			return fmt.Errorf("failed to marshal ClusterConfigSpec to ClusterVariable: %w", err)
		}
		cluster.Spec.Topology.Variables[clusterConfigVariableIndex] = *clusterConfigVariable

	}

	// Set defaults from 'workerConfig' variable from spec.topology.variables
	workerConfigVariable, workerConfigVariableIndex := variables.GetClusterVariableByName(
		workerconfig.MetaVariableName,
		cluster.Spec.Topology.Variables,
	)
	if workerConfigVariable != nil {
		workerConfigSpec := &v1alpha1.NodeConfigSpec{}
		err := workerConfigSpec.FromClusterVariable(workerConfigVariable)
		if err != nil {
			return fmt.Errorf("failed to unmarshal NodeConfigSpec from WorkerVariable: %w", err)
		}
		errs := defaultWorkerConfig(workerConfigSpec)
		if len(errs) > 1 {
			allErrs = append(allErrs, errs...)
		}
		workerConfigVariable, err = workerConfigSpec.ToClusterVariable(workerconfig.MetaVariableName)
		if err != nil {
			return fmt.Errorf("failed to marshal NodeConfigSpec to WorkerVariable: %w", err)
		}
		cluster.Spec.Topology.Variables[workerConfigVariableIndex] = *workerConfigVariable
	}

	// Set defaults for 'workerConfig' variable from spec.topology.workers.machineDeployments.variables.overrides
	if cluster.Spec.Topology.Workers != nil {
		for i, md := range cluster.Spec.Topology.Workers.MachineDeployments {
			mdWorkerConfigVariable, mdWorkerConfigVariableIndex := variables.GetMachineDeploymentVariableByName(
				workerconfig.MetaVariableName,
				md.Variables,
			)
			if mdWorkerConfigVariable != nil {
				workerConfigSpec := &v1alpha1.NodeConfigSpec{}
				err := workerConfigSpec.FromClusterVariable(mdWorkerConfigVariable)
				if err != nil {
					return fmt.Errorf("failed to unmarshal NodeConfigSpec from WorkerVariable: %w", err)
				}
				errs := defaultWorkerConfig(workerConfigSpec)
				if len(errs) > 1 {
					allErrs = append(allErrs, errs...)
				}
				mdWorkerConfigVariable, err = workerConfigSpec.ToClusterVariable(workerconfig.MetaVariableName)
				if err != nil {
					return fmt.Errorf("failed to marshal NodeConfigSpec to WorkerVariable: %w", err)
				}
				cluster.Spec.Topology.Workers.MachineDeployments[i].Variables.Overrides[mdWorkerConfigVariableIndex] =
					*mdWorkerConfigVariable
			}
		}
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(clusterv1.GroupVersion.WithKind("Cluster").GroupKind(), cluster.Name, allErrs)
	}

	return nil
}

func defaultClusterConfig(clusterConfig *v1alpha1.ClusterConfigSpec) field.ErrorList {
	// Set defaults for clusterConfig variable here.

	return nil
}

func defaultWorkerConfig(workerConfig *v1alpha1.NodeConfigSpec) field.ErrorList {
	// Set defaults for workerConfig variable here.

	return nil
}

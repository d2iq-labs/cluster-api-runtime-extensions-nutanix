// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nutanix

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	caaphv1 "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/external/sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/api/v1alpha1"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/k8s/client"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/config"
	lifecycleutils "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/lifecycle/utils"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/options"
)

const (
	defaultStorageHelmReleaseName      = "nutanix-csi-storage"
	defaultStorageHelmReleaseNamespace = "ntnx-system"

	defaultSnapshotHelmReleaseName      = "nutanix-csi-snapshot"
	defaultSnapshotHelmReleaseNamespace = "ntnx-system"

	//nolint:gosec // Does not contain hard coded credentials.
	defaultCredentialsSecretName = "nutanix-csi-credentials"
)

var defaultStorageClassParameters = map[string]string{
	"storageType":                                           "NutanixVolumes",
	"csi.storage.k8s.io/fstype":                             "xfs",
	"csi.storage.k8s.io/provisioner-secret-name":            defaultCredentialsSecretName,
	"csi.storage.k8s.io/provisioner-secret-namespace":       defaultStorageHelmReleaseNamespace,
	"csi.storage.k8s.io/node-publish-secret-name":           defaultCredentialsSecretName,
	"csi.storage.k8s.io/node-publish-secret-namespace":      defaultStorageHelmReleaseNamespace,
	"csi.storage.k8s.io/controller-expand-secret-name":      defaultCredentialsSecretName,
	"csi.storage.k8s.io/controller-expand-secret-namespace": defaultStorageHelmReleaseNamespace,
}

type NutanixCSIConfig struct {
	*options.GlobalOptions
	defaultValuesTemplateConfigMapName string
}

func (n *NutanixCSIConfig) AddFlags(prefix string, flags *pflag.FlagSet) {
	flags.StringVar(
		&n.defaultValuesTemplateConfigMapName,
		prefix+".default-values-template-configmap-name",
		"default-nutanix-csi-helm-values-template",
		"default values ConfigMap name",
	)
}

type NutanixCSI struct {
	client              ctrlclient.Client
	config              *NutanixCSIConfig
	helmChartInfoGetter *config.HelmChartGetter
}

func New(
	c ctrlclient.Client,
	cfg *NutanixCSIConfig,
	helmChartInfoGetter *config.HelmChartGetter,
) *NutanixCSI {
	return &NutanixCSI{
		client:              c,
		config:              cfg,
		helmChartInfoGetter: helmChartInfoGetter,
	}
}

func (n *NutanixCSI) Apply(
	ctx context.Context,
	provider v1alpha1.CSIProvider,
	defaultStorageConfig *v1alpha1.DefaultStorage,
	req *runtimehooksv1.AfterControlPlaneInitializedRequest,
) error {
	strategy := provider.Strategy
	switch strategy {
	case v1alpha1.AddonStrategyHelmAddon:
		err := n.handleHelmAddonApply(ctx, req)
		if err != nil {
			return err
		}
	case v1alpha1.AddonStrategyClusterResourceSet:
	default:
		return fmt.Errorf("stategy %s not implemented", strategy)
	}

	if provider.Credentials != nil {
		key := ctrlclient.ObjectKey{
			Name:      defaultCredentialsSecretName,
			Namespace: defaultStorageHelmReleaseNamespace,
		}
		err := lifecycleutils.CopySecretToRemoteCluster(
			ctx,
			n.client,
			provider.Credentials.Name,
			key,
			&req.Cluster,
		)
		if err != nil {
			return fmt.Errorf(
				"error creating credentials Secret for the Nutanix CSI driver: %w",
				err,
			)
		}
	}

	err := n.createStorageClasses(
		ctx,
		provider.StorageClassConfig,
		&req.Cluster,
		defaultStorageConfig,
	)
	if err != nil {
		return fmt.Errorf("error creating StorageClasses for the Nutanix CSI driver: %w", err)
	}

	return nil
}

func (n *NutanixCSI) handleHelmAddonApply(
	ctx context.Context,
	req *runtimehooksv1.AfterControlPlaneInitializedRequest,
) error {
	valuesTemplateConfigMap, err := lifecycleutils.RetrieveValuesTemplateConfigMap(ctx,
		n.client,
		n.config.defaultValuesTemplateConfigMapName,
		n.config.DefaultsNamespace())
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve nutanix csi installation values template ConfigMap for cluster: %w",
			err,
		)
	}
	values := valuesTemplateConfigMap.Data["values.yaml"]
	log := ctrl.LoggerFrom(ctx).WithValues(
		"cluster",
		ctrlclient.ObjectKeyFromObject(&req.Cluster),
	)
	helmChart, err := n.helmChartInfoGetter.For(ctx, log, config.NutanixStorageCSI)
	if err != nil {
		return fmt.Errorf("failed to get values for nutanix-csi-config %w", err)
	}

	hcp := &caaphv1.HelmChartProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: caaphv1.GroupVersion.String(),
			Kind:       "HelmChartProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Cluster.Namespace,
			Name:      "nutanix-csi-" + req.Cluster.Name,
		},
		Spec: caaphv1.HelmChartProxySpec{
			RepoURL:   helmChart.Repository,
			ChartName: helmChart.Name,
			ClusterSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{clusterv1.ClusterNameLabel: req.Cluster.Name},
			},
			ReleaseNamespace: defaultStorageHelmReleaseNamespace,
			ReleaseName:      defaultStorageHelmReleaseName,
			Version:          helmChart.Version,
			ValuesTemplate:   values,
		},
	}

	if err = controllerutil.SetOwnerReference(&req.Cluster, hcp, n.client.Scheme()); err != nil {
		return fmt.Errorf(
			"failed to set owner reference on nutanix-csi installation HelmChartProxy: %w",
			err,
		)
	}

	if err = client.ServerSideApply(ctx, n.client, hcp); err != nil {
		return fmt.Errorf("failed to apply nutanix-csi installation HelmChartProxy: %w", err)
	}

	snapshotHelmChart, err := n.helmChartInfoGetter.For(ctx, log, config.NutanixSnapshotCSI)
	if err != nil {
		return fmt.Errorf("failed to get values for nutanix-csi-config %w", err)
	}

	snapshotChart := &caaphv1.HelmChartProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: caaphv1.GroupVersion.String(),
			Kind:       "HelmChartProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Cluster.Namespace,
			Name:      "nutanix-csi-snapshot-" + req.Cluster.Name,
		},
		Spec: caaphv1.HelmChartProxySpec{
			RepoURL:   snapshotHelmChart.Repository,
			ChartName: snapshotHelmChart.Name,
			ClusterSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{clusterv1.ClusterNameLabel: req.Cluster.Name},
			},
			ReleaseNamespace: defaultSnapshotHelmReleaseNamespace,
			ReleaseName:      defaultSnapshotHelmReleaseName,
			Version:          snapshotHelmChart.Version,
		},
	}

	if err = client.ServerSideApply(ctx, n.client, snapshotChart); err != nil {
		return fmt.Errorf(
			"failed to apply nutanix-csi-snapshot installation HelmChartProxy: %w",
			err,
		)
	}

	return nil
}

func (n *NutanixCSI) createStorageClasses(
	ctx context.Context,
	configs []v1alpha1.StorageClassConfig,
	cluster *clusterv1.Cluster,
	defaultStorageConfig *v1alpha1.DefaultStorage,
) error {
	allStorageClasses := make([]runtime.Object, 0, len(configs))
	for _, config := range configs {
		setAsDefault := config.Name == defaultStorageConfig.StorageClassConfigName &&
			v1alpha1.CSIProviderNutanix == defaultStorageConfig.ProviderName
		allStorageClasses = append(allStorageClasses, lifecycleutils.CreateStorageClass(
			config,
			v1alpha1.NutanixProvisioner,
			setAsDefault,
			defaultStorageClassParameters,
		))
	}
	cm, err := lifecycleutils.CreateConfigMapForCRS(
		fmt.Sprintf("nutanix-storageclass-cm-%s", cluster.Name),
		n.config.DefaultsNamespace(),
		allStorageClasses...,
	)
	if err != nil {
		return err
	}
	err = client.ServerSideApply(ctx, n.client, cm)
	if err != nil {
		return err
	}
	return lifecycleutils.EnsureCRSForClusterFromObjects(
		ctx,
		"nutanix-storageclass-crs",
		n.client,
		cluster,
		cm,
	)
}

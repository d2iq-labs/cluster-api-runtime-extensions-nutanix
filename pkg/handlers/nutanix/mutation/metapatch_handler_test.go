// Copyright 2024 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mutation

import (
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/common/pkg/capi/clustertopology/handlers/mutation"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/clusterconfig"
	auditpolicytests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/auditpolicy/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/etcd"
	etcdtests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/etcd/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/extraapiservercertsans"
	extraapiservercertsanstests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/extraapiservercertsans/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/httpproxy"
	httpproxytests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/httpproxy/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/imageregistries"
	imageregistrycredentialstests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/imageregistries/credentials/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/kubernetesimagerepository"
	kubernetesimagerepositorytests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/kubernetesimagerepository/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/mirrors"
	globalimageregistrymirrortests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/generic/mutation/mirrors/tests"
	nutanixclusterconfig "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/nutanix/clusterconfig"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/nutanix/mutation/controlplaneendpoint"
	controlplaneendpointtests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/nutanix/mutation/controlplaneendpoint/tests"
	"github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/nutanix/mutation/prismcentralendpoint"
	prismcentralendpointtests "github.com/d2iq-labs/cluster-api-runtime-extensions-nutanix/pkg/handlers/nutanix/mutation/prismcentralendpoint/tests"
)

func metaPatchGeneratorFunc(mgr manager.Manager) func() mutation.GeneratePatches {
	return func() mutation.GeneratePatches {
		return MetaPatchHandler(mgr).(mutation.GeneratePatches)
	}
}

func workerPatchGeneratorFunc() func() mutation.GeneratePatches {
	return func() mutation.GeneratePatches {
		return MetaWorkerPatchHandler().(mutation.GeneratePatches)
	}
}

func TestGeneratePatches(t *testing.T) {
	t.Parallel()

	mgr := testEnv.Manager

	controlplaneendpointtests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		nutanixclusterconfig.NutanixVariableName,
		controlplaneendpoint.VariableName,
	)

	prismcentralendpointtests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		nutanixclusterconfig.NutanixVariableName,
		prismcentralendpoint.VariableName,
	)

	auditpolicytests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
	)

	httpproxytests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		httpproxy.VariableName,
	)

	etcdtests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		etcd.VariableName,
	)

	extraapiservercertsanstests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		extraapiservercertsans.VariableName,
	)

	kubernetesimagerepositorytests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		clusterconfig.MetaVariableName,
		kubernetesimagerepository.VariableName,
	)

	imageregistrycredentialstests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		mgr.GetClient(),
		clusterconfig.MetaVariableName,
		imageregistries.VariableName,
	)

	globalimageregistrymirrortests.TestGeneratePatches(
		t,
		metaPatchGeneratorFunc(mgr),
		mgr.GetClient(),
		clusterconfig.MetaVariableName,
		mirrors.GlobalMirrorVariableName,
	)
}

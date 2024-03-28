# Copyright 2023 D2iQ, Inc. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

.PHONY: clusterctl.init
clusterctl.init:
	env CLUSTER_TOPOLOGY=true \
	    EXP_RUNTIME_SDK=true \
	    EXP_CLUSTER_RESOURCE_SET=true \
	    EXP_MACHINE_POOL=true \
	    AWS_B64ENCODED_CREDENTIALS=$$(clusterawsadm bootstrap credentials encode-as-profile) \
	    clusterctl init \
	      --kubeconfig=$(KIND_KUBECONFIG) \
	      --infrastructure docker,aws,nutanix \
	      --addon helm \
	      --wait-providers

.PHONY: clusterctl.delete
clusterctl.delete:
	clusterctl delete --kubeconfig=$(KIND_KUBECONFIG) --all

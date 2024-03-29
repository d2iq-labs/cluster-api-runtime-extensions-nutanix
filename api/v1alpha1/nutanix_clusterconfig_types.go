// Copyright 2024 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NutanixSpec defines the desired state of NutanixCluster
type NutanixSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// host can be either DNS name or ip address
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// Nutanix Prism Central endpoint configuration.
	PrismCentralEndpoint NutanixPrismCentralEndpointSpec `json:"prismCentral"`
}

func (NutanixSpec) VariableSchema() clusterv1.VariableSchema {
	return clusterv1.VariableSchema{
		OpenAPIV3Schema: clusterv1.JSONSchemaProps{
			Description: "Nutanix cluster configuration",
			Type:        "object",
			Properties: map[string]clusterv1.JSONSchemaProps{
				"prismCentralEndpoint": NutanixPrismCentralEndpointSpec{}.VariableSchema().OpenAPIV3Schema,
				"controlPlaneEndpoint": ControlPlaneEndpointSpec{}.VariableSchema().OpenAPIV3Schema,
			},
		},
	}
}

type NutanixPrismCentralEndpointSpec struct {
	// address is the endpoint address (DNS name or IP address) of the Nutanix Prism Central
	Address string `json:"address"`

	// port is the port number to access the Nutanix Prism Central
	Port int32 `json:"port"`

	// use insecure connection to Prism Central endpoint
	// +optional
	Insecure bool `json:"insecure"`

	// A reference to the ConfigMap containing a PEM encoded x509 cert for the RootCA that was used to create the certificate
	// for a Prism Central that uses certificates that were issued by a non-publicly trusted RootCA. The trust
	// bundle is added to the cert pool used to authenticate the TLS connection to the Prism Central.
	// +optional
	AdditionalTrustBundle *corev1.LocalObjectReference `json:"additionalTrustBundle,omitempty"`

	// A reference to the Secret for credential information for the target Prism Central instance
	Credentials corev1.LocalObjectReference `json:"credentials"`
}

func (NutanixPrismCentralEndpointSpec) VariableSchema() clusterv1.VariableSchema {
	return clusterv1.VariableSchema{
		OpenAPIV3Schema: clusterv1.JSONSchemaProps{
			Description: "Nutanix Prism Central endpoint configuration",
			Type:        "object",
			Properties: map[string]clusterv1.JSONSchemaProps{
				"address": {
					Description: "the endpoint address (DNS name or IP address) of the Nutanix Prism Central",
					Type:        "string",
				},
				"port": {
					Description: "The port number to access the Nutanix Prism Central",
					Type:        "integer",
				},
				"insecure": {
					Description: "Use insecure connection to Prism Central endpoint",
					Type:        "boolean",
				},
				"additionalTrustBundle": {
					Description: "A reference to the ConfigMap containing a PEM encoded x509 cert for the RootCA " +
						"that was used to create the certificate for a Prism Central that uses certificates " +
						"that were issued by a non-publicly trusted RootCA." +
						"The trust bundle is added to the cert pool used to authenticate the TLS connection " +
						"to the Prism Central.",
					Type: "object",
					Properties: map[string]clusterv1.JSONSchemaProps{
						"name": {
							Description: "The name of the ConfigMap",
							Type:        "string",
						},
					},
					Required: []string{"name"},
				},
				"credentials": {
					Description: "A reference to the Secret for credential information" +
						"for the target Prism Central instance",
					Type: "object",
					Properties: map[string]clusterv1.JSONSchemaProps{
						"name": {
							Description: "The name of the Secret",
							Type:        "string",
						},
					},
					Required: []string{"name"},
				},
			},
			Required: []string{"address", "port", "credentials"},
		},
	}
}

// Copyright 2023 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package webhooks

import "github.com/spf13/pflag"

type Options struct {
	WebhookPort    int
	WebhookCertDir string
}

func NewOptions() *Options {
	return &Options{}
}

func (s *Options) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&s.WebhookPort, "webhook-port", s.WebhookPort, "Webhook server port")

	fs.StringVar(
		&s.WebhookCertDir,
		"webhook-cert-dir",
		s.WebhookCertDir,
		"Webhook server cert dir",
	)
}

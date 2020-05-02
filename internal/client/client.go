// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

// New create new VMWare client
func New(vmURL string, vmUsername string, vmPassword string, ValidateSSL bool) (*govmomi.Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // Parse URL from string
	url, err := soap.ParseURL(vmURL)
	if err != nil {
		return nil, err
	}

	// Override username and/or password as required
	setCredentials(url, vmUsername, vmPassword)

	// Connect and log in to ESX/i or vCenter
	return govmomi.NewClient(ctx, url, !ValidateSSL)
}

func setCredentials(u *url.URL, un string, pw string) {
	// Override username if provided
	if un != "" {
		var password string
		var ok bool
		if u.User != nil {
			password, ok = u.User.Password()
		}
		if ok {
			u.User = url.UserPassword(un, password)
		} else {
			u.User = url.User(un)
		}
	}

	// Override password if provided
	if pw != "" {
		var username string
		if u.User != nil {
			username = u.User.Username()
		}
		u.User = url.UserPassword(username, pw)
	}
}

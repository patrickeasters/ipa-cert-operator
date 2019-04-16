// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package ipa

import (
	"encoding/json"
)

// HostRecord encapsulates host data returned from ipa host commands
type HostRecord struct {
	Dn          string    `json:"dn"`
	Fqdn        IpaString `json:"fqdn"`
	Description IpaString `json:"description"`
	HasKeytab   bool      `json:"has_keytab"`
	HasPassword bool      `json:"has_password"`
}

// Fetch host details by calling the FreeIPA host-show method
func (c *Client) HostShow(fqdn string) (*HostRecord, error) {

	options := map[string]interface{}{
		"all": true}

	res, err := c.rpc("host_show", []string{fqdn}, options)

	if err != nil {
		return nil, err
	}

	var hostRec HostRecord
	err = json.Unmarshal(res.Result.Data, &hostRec)
	if err != nil {
		return nil, err
	}

	return &hostRec, nil
}

// HostAdd creates a new host record
func (c *Client) HostAdd(fqdn, description string, skip_dns bool) (*HostRecord, error) {
	var options = map[string]interface{}{
		"force": skip_dns}

	if len(description) > 0 {
		options["description"] = description
	}

	res, err := c.rpc("host_add", []string{fqdn}, options)
	if err != nil {
		return nil, err
	}

	var hostRec HostRecord
	err = json.Unmarshal(res.Result.Data, &hostRec)
	if err != nil {
		return nil, err
	}

	return &hostRec, nil
}

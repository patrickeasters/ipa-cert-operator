// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package ipa

import (
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"
)

// CertRecord encapsulates cert data returned from ipa cert commands
type CertRecord struct {
	CaCn        string `json:"cacn"`
	Certificate string `json:"certificate"`
	Issuer      string `json:"issuer"`
	Subject     string `json:"subject"`
	StartDate   string `json:"valid_not_before"`
	EndDate     string `json:"valid_not_after"`
	Serial      string `json:"serial_number_hex"`
	Revoked     bool   `json:"revoked"`
}

const certTimeFormat = "Mon Jan 02 15:04:05 2006 MST"

// Convert DER cert to PEM format
func (c *CertRecord) CertPem() (string, error) {
	decoded, err := b64.StdEncoding.DecodeString(c.Certificate)
	if err != nil {
		return "", fmt.Errorf("Unable to decode certificate: %s", err)
	}

	cpem := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: decoded}))
	return cpem, nil

}

// Get expiration date from cert record
func (cert *CertRecord) Expiry() (*time.Time, error) {
	t, err := time.Parse(certTimeFormat, cert.EndDate)
	return &t, err
}

// Get issued date from cert record
func (cert *CertRecord) Issued() (*time.Time, error) {
	t, err := time.Parse(certTimeFormat, cert.StartDate)
	return &t, err
}

// Checks whether cert is expired or revoked
func (cert *CertRecord) Valid() (bool, error) {
	exp, err := cert.Expiry()
	if err != nil {
		return true, err
	}

	if exp.Before(time.Now()) || cert.Revoked {
		return false, nil
	}

	return true, nil
}

// Fetch cert details by call the FreeIPA cert-show method
func (c *Client) CertShow(serial string) (*CertRecord, error) {
	options := map[string]interface{}{
		"all": true}

	res, err := c.rpc("cert_show", []string{serial}, options)

	if err != nil {
		return nil, err
	}

	var certRec CertRecord
	err = json.Unmarshal(res.Result.Data, &certRec)
	if err != nil {
		return nil, err
	}

	return &certRec, nil
}

func (c *Client) CertRevoke(serial string) error {
	options := map[string]interface{}{}

	_, err := c.rpc("cert_revoke", []string{serial}, options)

	return err
}

func (c *Client) CertRequest(principal, csr, profile string) (*CertRecord, error) {
	options := map[string]interface{}{
		"all":        true,
		"principal":  principal,
		"profile_id": profile}

	res, err := c.rpc("cert_request", []string{csr}, options)

	if err != nil {
		return nil, err
	}

	var certRec CertRecord
	err = json.Unmarshal(res.Result.Data, &certRec)
	if err != nil {
		return nil, err
	}

	return &certRec, nil
}

package fawry

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiPath        = "/ECommerceWeb/Fawry/payments/"
	baseURL        = "https://www.atfawry.com"
	sandboxBaseURL = "https://atfawry.fawrystaging.com"
)

// Client Struct
type Client struct {
	IsSandbox      bool
	FawrySecureKey string
	httpClient     *http.Client
}

// NewClient returns a new Fawry API client.
func NewClient(isSandbox bool, secureKey string) *Client {
	return &Client{
		IsSandbox:      isSandbox,
		FawrySecureKey: secureKey,
		httpClient:     &http.Client{},
	}
}

func (fc Client) getURL() string {
	if fc.IsSandbox {
		return sandboxBaseURL + apiPath
	}
	return baseURL + apiPath
}

func (fc Client) getSignature(inputs []string) string {
	sum := sha256.Sum256([]byte(strings.Join(inputs[:], ",")))
	return hex.EncodeToString(sum[:])
}

// ChargeRequest could be used to charge the customer with different payment methods.
//
//	It also might be used to create a reference number to be paid at Fawry's outlets or
//	it can be used to direct debit the customer card using card token.
func (fc Client) ChargeRequest(charge Charge) (*http.Response, error) {
	err := charge.Validate()
	if err != nil {
		return nil, err
	}

	url := fc.getURL() + "charge"

	signatureArray := []string{charge.MerchantCode,
		charge.MerchantRefNum, charge.CustomerProfileID,
		charge.PaymentMethod, charge.Amount, charge.CardToken, fc.FawrySecureKey}

	body := struct {
		Charge
		Signature string `json:"signature"`
	}{Charge: charge,
		Signature: fc.getSignature(signatureArray)}

	return fc.sendRequest(http.MethodPost, url, body)
}

// RefundRequest  can refund the payment again to the customer
func (fc Client) RefundRequest(refund Refund) (*http.Response, error) {
	err := refund.Validate()
	if err != nil {
		return nil, err
	}

	url := fc.getURL() + "refund"

	signatureArray := []string{refund.MerchantCode,
		refund.ReferenceNumber, refund.RefundAmount,
		refund.Reason, fc.FawrySecureKey}

	body := struct {
		Refund
		Signature string `json:"signature"`
	}{
		Refund:    refund,
		Signature: fc.getSignature(signatureArray),
	}

	return fc.sendRequest(http.MethodPost, url, body)

}

// StatusRequest can use Get Payment Status Service to retrieve the payment status for the charge request
func (fc Client) StatusRequest(status Status) (*http.Response, error) {
	err := status.Validate()
	if err != nil {
		return nil, err
	}

	signatureArray := []string{status.MerchantCode, status.MerchantRefNum, fc.FawrySecureKey}

	url := fc.getURL() + "status" + fmt.Sprintf("?merchantCode=%s&merchantRefNumber=%s&signature=%s",
		status.MerchantCode,
		status.MerchantRefNum,
		fc.getSignature(signatureArray))

	return fc.sendRequest(http.MethodGet, url, nil)

}

// sendRequest is a helper function to send HTTP requests
func (fc Client) sendRequest(method, url string, body interface{}) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	return resp, nil
}

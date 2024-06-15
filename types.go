package fawry

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// ChargeItem Struct
type ChargeItem struct {
	ItemID      string `json:"itemId"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Quantity    int    `json:"quantity"`
}

// Validate func for ChargeItem struct
func (chargeItem ChargeItem) Validate() error {
	return validation.ValidateStruct(&chargeItem,
		validation.Field(&chargeItem.ItemID, validation.Required),
		validation.Field(&chargeItem.Description, validation.Required),
		validation.Field(&chargeItem.Price, validation.Required, validation.Match(regexp.MustCompile(`^\d+\.\d\d$`))),
		validation.Field(&chargeItem.Quantity, validation.Required),
	)
}

// Charge Struct
type Charge struct {
	MerchantCode      string       `json:"merchantCode"`
	MerchantRefNum    string       `json:"merchantRefNum"`
	CustomerProfileID string       `json:"customerProfileId"`
	Amount            string       `json:"amount"`
	Description       string       `json:"description"`
	CustomerMobile    string       `json:"customerMobile"`
	ChargeItems       []ChargeItem `json:"chargeItems"`
	CurrencyCode      string       `json:"currencyCode"`
	CardToken         string       `json:"cardToken"`
	CustomerEmail     string       `json:"customerEmail"`
	PaymentMethod     string       `json:"paymentMethod"`
	PaymentExpiry     int          `json:"paymentExpiry"`
}

// Customized validation func to check card if payment method is card
func cardTokenRequired(paymentMethod string) validation.RuleFunc {
	return func(value interface{}) error {
		cardToken, _ := value.(string)
		if paymentMethod == "CARD" && len(cardToken) == 0 {
			return errors.New("card token is required when payment method is card")
		}
		return nil
	}
}

// Validate func for Charge Struct
func (charge Charge) Validate() error {
	return validation.ValidateStruct(&charge,
		validation.Field(&charge.MerchantCode, validation.Required),
		validation.Field(&charge.MerchantRefNum, validation.Required),
		validation.Field(&charge.CustomerProfileID, validation.Required),
		validation.Field(&charge.Amount, validation.Required, validation.Match(regexp.MustCompile(`^\d+\.\d\d$`))),
		validation.Field(&charge.Description, validation.Required),
		validation.Field(&charge.CustomerMobile, validation.Required),
		validation.Field(&charge.ChargeItems, validation.Required),
		validation.Field(&charge.CurrencyCode),
		validation.Field(&charge.CardToken),
		validation.Field(&charge.CustomerEmail, is.Email),
		validation.Field(&charge.PaymentMethod),
		validation.Field(&charge.PaymentExpiry),
		validation.Field(&charge.CardToken, validation.By(cardTokenRequired(charge.PaymentMethod))),
	)
}

// Refund Struct
type Refund struct {
	MerchantCode    string `json:"merchantCode"`
	ReferenceNumber string `json:"referenceNumber"`
	RefundAmount    string `json:"refundAmount"`
	Reason          string `json:"reason"`
}

// Validate func for Refund Struct
func (refund Refund) Validate() error {
	return validation.ValidateStruct(&refund,
		validation.Field(&refund.MerchantCode, validation.Required),
		validation.Field(&refund.ReferenceNumber, validation.Required),
		validation.Field(&refund.RefundAmount, validation.Required, validation.Match(regexp.MustCompile(`^\d+\.\d\d$`))),
		validation.Field(&refund.Reason, validation.Required),
	)
}

// Status Struct
type Status struct {
	MerchantCode   string `json:"merchantCode"`
	MerchantRefNum string `json:"merchantRefNum"`
}

// Validate func for Status Struct
func (status Status) Validate() error {
	return validation.ValidateStruct(&status,
		validation.Field(&status.MerchantCode, validation.Required),
		validation.Field(&status.MerchantRefNum, validation.Required),
	)
}

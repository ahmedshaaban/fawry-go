package fawry

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
)

type fields struct {
	IsSandbox      bool
	FawrySecureKey string
}

func TestClient_getURL(t *testing.T) {

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "returns production link if isSandbox false",
			fields: fields{
				IsSandbox:      false,
				FawrySecureKey: "test",
			},
			want: "https://www.atfawry.com/ECommerceWeb/Fawry/payments/",
		},
		{
			name: "returns sandbox link if isSandbox true",
			fields: fields{
				IsSandbox:      true,
				FawrySecureKey: "test",
			},
			want: "https://atfawry.fawrystaging.com/ECommerceWeb/Fawry/payments/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewClient(tt.fields.IsSandbox, tt.fields.FawrySecureKey)
			if got := fc.getURL(); got != tt.want {
				t.Errorf("Client.getURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_getSignature(t *testing.T) {
	type args struct {
		inputs []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "correctly hash values",
			fields: fields{
				IsSandbox:      true,
				FawrySecureKey: "testSecureKey",
			},
			args: args{
				inputs: []string{
					"test1", "test2", "test3",
				},
			},
			want: "b20c54f469f62a14076062e93962329c1b13353ba19e48088b87559a0955465a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewClient(tt.fields.IsSandbox, tt.fields.FawrySecureKey)
			if got := fc.getSignature(tt.args.inputs); got != tt.want {
				t.Errorf("Client.getSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ChargeRequest(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	mockedRespBody := map[string]interface{}{
		"type":              "ChargeResponse",
		"referenceNumber":   "100163201",
		"merchantRefNumber": "9990064204",
		"expirationTime":    1516554874077,
		"statusCode":        200,
		"statusDescription": "Operation done successfully"}

	gock.New("https://www.atfawry.com/ECommerceWeb/Fawry/payments").
		Post("/charge").
		Reply(200).
		JSON(mockedRespBody)

	gock.New("https://atfawry.fawrystaging.com/ECommerceWeb/Fawry/payments").
		Post("/charge").
		Reply(200).
		JSON(mockedRespBody)

	type args struct {
		charge Charge
	}

	input := args{
		charge: Charge{
			MerchantCode:      "is0N+YQzlE4=",
			MerchantRefNum:    "9990064204",
			CustomerProfileID: "9990064204",
			CustomerMobile:    "01000000200",
			CustomerEmail:     "77@test.com",
			PaymentMethod:     "PAYATFAWRY",
			Amount:            "20.10",
			CurrencyCode:      "EGP",
			Description:       "the charge request description",
			PaymentExpiry:     1516554874077,
			ChargeItems: []ChargeItem{
				{ItemID: "897fa8e81be26df25db592e81c31c",
					Description: "asdasd",
					Price:       "15.20",
					Quantity:    1}},
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "test sandbox charge request",
			fields: fields{
				IsSandbox:      true,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
		{
			name: "test production charge request",
			fields: fields{
				IsSandbox:      false,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewClient(tt.fields.IsSandbox, tt.fields.FawrySecureKey)

			got, _ := fc.ChargeRequest(tt.args.charge)

			if got.Status != tt.want.Status {
				t.Errorf("Client.ChargeRequest() got = %v, expected %v", got.Status, tt.want.Status)
			}
		})
	}
}

func TestClient_RefundRequest(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	mockedRespBody := map[string]interface{}{
		"type":              "ResponseDataModel",
		"statusCode":        200,
		"statusDescription": "Operation done successfully",
	}

	gock.New("https://www.atfawry.com/ECommerceWeb/Fawry/payments").
		Post("/refund").
		Reply(200).
		JSON(mockedRespBody)

	gock.New("https://atfawry.fawrystaging.com/ECommerceWeb/Fawry/payments").
		Post("/refund").
		Reply(200).
		JSON(mockedRespBody)

	type args struct {
		refund Refund
	}

	input := args{
		refund: Refund{
			MerchantCode:    "1013969",
			ReferenceNumber: "322818",
			RefundAmount:    "100.00",
			Reason:          "Bad Quality ",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "test sandbox refund request",
			fields: fields{
				IsSandbox:      true,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
		{
			name: "test production refund request",
			fields: fields{
				IsSandbox:      false,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewClient(tt.fields.IsSandbox, tt.fields.FawrySecureKey)
			got, _ := fc.RefundRequest(tt.args.refund)
			if got.Status != tt.want.Status {
				t.Errorf("Client.RefundRequest() got = %v, expected %v", got.Status, tt.want.Status)
			}
		})
	}
}

func TestClient_StatusRequest(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	mockedRespBody := map[string]interface{}{
		"type":              "PaymentStatusResponse",
		"referenceNumber":   "100162801",
		"merchantRefNumber": "99900642041",
		"paymentAmount":     20,
		"paymentDate":       1514747471138,
		"expirationTime":    151654747115,
		"paymentStatus":     "PAID",
		"paymentMethod":     "PAYATFAWRY",
		"statusCode":        200,
		"statusDescription": "Operation done successfully"}

	gock.New("https://www.atfawry.com/ECommerceWeb/Fawry/payments").
		Get("/status").
		Reply(200).
		JSON(mockedRespBody)

	gock.New("https://atfawry.fawrystaging.com/ECommerceWeb/Fawry/payments").
		Get("/status").
		Reply(200).
		JSON(mockedRespBody)

	type args struct {
		status Status
	}

	input := args{
		status: Status{
			MerchantCode:   "is0N+YQzlE4=",
			MerchantRefNum: "99900642041",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "test sandbox status request",
			fields: fields{
				IsSandbox:      true,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
		{
			name: "test production status request",
			fields: fields{
				IsSandbox:      false,
				FawrySecureKey: "testSecuredKey",
			},
			args: input,
			want: &http.Response{
				Status: "200 OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewClient(tt.fields.IsSandbox, tt.fields.FawrySecureKey)
			got, _ := fc.StatusRequest(tt.args.status)
			if got.Status != tt.want.Status {
				t.Errorf("Client.StatusRequest() got = %v, expected %v", got.Status, tt.want.Status)
			}
		})
	}
}

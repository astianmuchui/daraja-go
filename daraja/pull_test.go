package daraja

import (
	"encoding/json"
	"testing"
)

func TestPulledTransactionListFlattensNestedArrays(t *testing.T) {
	body := `{"ResponseCode":"1000","ResponseMessage":"Success","Response":[[{"transactionId":"yzlCSRpyQd","trxDate":"2019-07-31T20:35:21.000Z","msisdn":70000000,"sender":"FirstName","transactiontype":"c2b-pay-bill-debit","billreference":"TestAPI","amount":10,"organizationname":"Organization Name"}]]}`

	var payload QueryPullTransactionsResponsePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if payload.ResponseCode != "1000" {
		t.Errorf("ResponseCode = %q, want %q", payload.ResponseCode, "1000")
	}

	if len(payload.Response) != 1 {
		t.Fatalf("len(Response) = %d, want 1", len(payload.Response))
	}

	tx := payload.Response[0]

	if tx.TransactionID != "yzlCSRpyQd" {
		t.Errorf("TransactionID = %q, want %q", tx.TransactionID, "yzlCSRpyQd")
	}
	if tx.Amount.String() != "10" {
		t.Errorf("Amount = %q, want %q", tx.Amount.String(), "10")
	}
	if tx.MSISDN.String() != "70000000" {
		t.Errorf("MSISDN = %q, want %q", tx.MSISDN.String(), "70000000")
	}
	if tx.BillReference != "TestAPI" {
		t.Errorf("BillReference = %q, want %q", tx.BillReference, "TestAPI")
	}
	if tx.TransactionType != "c2b-pay-bill-debit" {
		t.Errorf("TransactionType = %q, want %q", tx.TransactionType, "c2b-pay-bill-debit")
	}
}

func TestPulledTransactionListAcceptsFlatArray(t *testing.T) {
	body := `{"Response":[{"transactionId":"A1","amount":"250.00","msisdn":"254700000000"}]}`

	var payload QueryPullTransactionsResponsePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(payload.Response) != 1 {
		t.Fatalf("len(Response) = %d, want 1", len(payload.Response))
	}
	if payload.Response[0].TransactionID != "A1" {
		t.Errorf("TransactionID = %q, want %q", payload.Response[0].TransactionID, "A1")
	}
	if payload.Response[0].Amount.String() != "250.00" {
		t.Errorf("Amount = %q, want %q", payload.Response[0].Amount.String(), "250.00")
	}
}

func TestPulledTransactionListAcceptsNullAndEmpty(t *testing.T) {
	for _, body := range []string{`{"Response":null}`, `{"Response":[]}`, `{"Response":[[]]}`, `{}`} {
		var payload QueryPullTransactionsResponsePayload
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			t.Fatalf("unmarshal %s: %v", body, err)
		}
		if len(payload.Response) != 0 {
			t.Errorf("len(Response) = %d for %s, want 0", len(payload.Response), body)
		}
	}
}

func TestPulledTransactionListRejectsGarbage(t *testing.T) {
	var payload QueryPullTransactionsResponsePayload
	if err := json.Unmarshal([]byte(`{"Response":"not a list"}`), &payload); err == nil {
		t.Error("unmarshal of a string Response succeeded, want an error")
	}
}

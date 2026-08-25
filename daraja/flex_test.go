package daraja

import (
	"encoding/json"
	"testing"
)


const accountBalanceResult = `{"Result":{"ResultType":0,"ResultCode":0,"ResultDesc":"The service request is processed successfully.","OriginatorConversationID":"8d0a-4cf1-a3f3-ebad341607ad16074","ConversationID":"AG_20260825_010020480ucukatfpxp2","TransactionID":"UHP0000000","ResultParameters":{"ResultParameter":[{"Key":"ActionType","Value":"AccountBalance"},{"Key":"AccountBalance","Value":"Working Account|KES|0.00|0.00|0.00|0.00&Utility Account|KES|0.00|0.00|0.00|0.00&Charges Paid Account|KES|0.00|0.00|0.00|0.00&Merchant Account|KES|0.00|0.00|0.00|0.00&Organization Settlement Account|KES|0.00|0.00|0.00|0.00"},{"Key":"BOCompletedTime","Value":20260825140947}]},"ReferenceData":{"ReferenceItem":{"Key":"QueueTimeoutURL","Value":"https://internalapi.safaricom.co.ke/mpesa/abresults/v1/submit"}}}}`

func TestCallbackPayloadDecodesAccountBalanceResult(t *testing.T) {
	var p CallbackPayload
	if err := json.Unmarshal([]byte(accountBalanceResult), &p); err != nil {
		t.Fatalf("account balance result did not decode: %v", err)
	}

	if got := p.Result.ResultCode.String(); got != "0" {
		t.Errorf("ResultCode = %q, want \"0\"", got)
	}
	if got := p.Result.ResultType.String(); got != "0" {
		t.Errorf("ResultType = %q, want \"0\"", got)
	}
	if p.Result.OriginatorConversationID != "8d0a-4cf1-a3f3-ebad341607ad16074" {
		t.Errorf("OriginatorConversationID = %q", p.Result.OriginatorConversationID)
	}

	params := p.Result.ResultParameters.ResultParameter
	if len(params) != 3 {
		t.Fatalf("got %d result parameters, want 3", len(params))
	}
	if got := params.Get("ActionType"); got != "AccountBalance" {
		t.Errorf("ActionType = %q", got)
	}

	// The 14-digit timestamp must survive verbatim: a float64 round-trip would
	// render it as 2.0260825140947e+13.
	if got := params.Get("BOCompletedTime"); got != "20260825140947" {
		t.Errorf("BOCompletedTime = %q, want the exact digits", got)
	}

	balances := params.Get("AccountBalance")
	if balances == "" {
		t.Fatal("AccountBalance parameter is empty")
	}

	// ReferenceItem arrived as a single object, not an array.
	if len(p.Result.ReferenceData.ReferenceItem) != 1 {
		t.Fatalf("got %d reference items, want 1", len(p.Result.ReferenceData.ReferenceItem))
	}
	if got := p.Result.ReferenceData.ReferenceItem.Get("QueueTimeoutURL"); got == "" {
		t.Error("QueueTimeoutURL reference item is empty")
	}
}

// The same fields are quoted on other callbacks, and an array-shaped
// ReferenceItem is what reversal results send.
func TestCallbackPayloadDecodesQuotedAndArrayShapes(t *testing.T) {
	const body = `{"Result":{"ResultType":"0","ResultCode":"21","ResultDesc":"Failed","OriginatorConversationID":"abc","ConversationID":"AG_1","TransactionID":"T1","ResultParameters":{"ResultParameter":[{"Key":"Amount","Value":"250"}]},"ReferenceData":{"ReferenceItem":[{"Key":"QueueTimeoutURL","Value":"https://example.test"},{"Key":"Occasion","Value":"test"}]}}}`

	var p CallbackPayload
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatalf("quoted/array shape did not decode: %v", err)
	}
	if got := p.Result.ResultCode.String(); got != "21" {
		t.Errorf("ResultCode = %q, want \"21\"", got)
	}
	if len(p.Result.ReferenceData.ReferenceItem) != 2 {
		t.Errorf("got %d reference items, want 2", len(p.Result.ReferenceData.ReferenceItem))
	}
	if got := p.Result.ResultParameters.ResultParameter.Get("Amount"); got != "250" {
		t.Errorf("Amount = %q", got)
	}
}

// A single ResultParameter object rather than an array.
func TestResultParameterAcceptsLoneObject(t *testing.T) {
	const body = `{"Result":{"ResultCode":0,"ResultParameters":{"ResultParameter":{"Key":"Amount","Value":100}}}}`

	var p CallbackPayload
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatalf("lone result parameter did not decode: %v", err)
	}
	if got := p.Result.ResultParameters.ResultParameter.Get("Amount"); got != "100" {
		t.Errorf("Amount = %q, want \"100\"", got)
	}
}

func TestFlexStringEdgeCases(t *testing.T) {
	for _, tc := range []struct{ name, in, want string }{
		{"quoted", `"abc"`, "abc"},
		{"integer", `42`, "42"},
		{"large integer", `20260825140947`, "20260825140947"},
		{"decimal", `1.5`, "1.5"},
		{"boolean", `true`, "true"},
		{"null", `null`, ""},
		{"empty string", `""`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var f FlexString
			if err := json.Unmarshal([]byte(tc.in), &f); err != nil {
				t.Fatalf("unmarshal %s: %v", tc.in, err)
			}
			if f.String() != tc.want {
				t.Errorf("got %q, want %q", f.String(), tc.want)
			}
		})
	}
}

// Decoding then re-encoding must produce valid JSON, so a callback can be
// logged or stored in a jsonb column after a round trip.
func TestFlexStringMarshalsAsString(t *testing.T) {
	var p CallbackPayload
	if err := json.Unmarshal([]byte(accountBalanceResult), &p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !json.Valid(out) {
		t.Error("re-encoded callback is not valid JSON")
	}

	var again CallbackPayload
	if err := json.Unmarshal(out, &again); err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if again.Result.ResultCode != p.Result.ResultCode {
		t.Errorf("ResultCode changed across a round trip: %q -> %q", p.Result.ResultCode, again.Result.ResultCode)
	}
}

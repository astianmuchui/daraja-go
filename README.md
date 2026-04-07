# daraja-go

daraja-go is a Go library for interacting with the Safaricom Daraja API. It covers the full range of M-Pesa operations including STK push, business payments, transaction queries, reversals, account balance checks, and pull transaction history. The library has no external dependencies and handles OAuth2 token management automatically.

**Author:** Sebastian Muchui
**License:** MIT
**Module:** `github.com/astianmuchui/daraja-go`
**Minimum Go version:** 1.24

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Environment Selection](#environment-selection)
- [Authentication](#authentication)
- [API Methods](#api-methods)
  - [STK Push (Lipa na M-Pesa Online)](#stk-push-lipa-na-m-pesa-online)
  - [Business to Customer (B2C)](#business-to-customer-b2c)
  - [Business to Business (B2B)](#business-to-business-b2b)
  - [Register C2B URLs](#register-c2b-urls)
  - [Query Transaction Status](#query-transaction-status)
  - [Reverse a Transaction](#reverse-a-transaction)
  - [Query Account Balance](#query-account-balance)
  - [Query Online Transaction Status](#query-online-transaction-status)
  - [Pull Transaction History](#pull-transaction-history)
  - [C2B Validation](#c2b-validation)
- [Error Handling](#error-handling)
- [Types Reference](#types-reference)
- [Constants](#constants)

---

## Installation

```bash
go get github.com/astianmuchui/daraja-go
```

---

## Configuration

Before making any API calls, set the package-level credential variables. These are required for authentication and for building request payloads.

```go
import "github.com/astianmuchui/daraja-go/daraja"

daraja.CONSUMER_KEY    = "your_consumer_key"
daraja.CONSUMER_SECRET = "your_consumer_secret"
daraja.SHORTCODE       = "your_shortcode"
daraja.PASSKEY         = "your_passkey"
daraja.ACCOUNT_TYPE    = "your_account_type"
```

These values are obtained from the [Safaricom Developer Portal](https://developer.safaricom.co.ke/) when you register your application.

---

## Environment Selection

The library defaults to the sandbox environment on startup. Switch to production before making live API calls.

```go
// Use sandbox (default, no call needed)
daraja.Production(false)

// Use production
daraja.Production(true)

// Convenience alias for production
daraja.SetProductionMode()
```

Changing the environment reinitializes all API endpoint URLs. Sandbox base URL is `https://sandbox.safaricom.co.ke` and production is `https://api.safaricom.co.ke`.

---

## Authentication

Authentication is handled automatically. Every API method checks whether the current token has expired before making a request and re-authorizes if necessary. You rarely need to call `Authorize` directly.

```go
client := &daraja.Daraja{}

// Manual authorization (optional)
ok, errs := client.Authorize()
if !ok {
    for _, e := range errs {
        log.Println(e)
    }
}

// Check if the current token is still valid
if client.IsAuthorized() {
    // Token is still active
}

// Asynchronous authorization using channels
statusCh := make(chan bool)
errsCh   := make(chan []error)
client.RetryAuth(statusCh, errsCh)
ok   = <-statusCh
errs = <-errsCh
```

`Authorize` uses HTTP Basic Auth with `CONSUMER_KEY` and `CONSUMER_SECRET` to fetch an OAuth2 bearer token. The token expiry is tracked internally using the `expires_in` value returned by the API.

---

## API Methods

All API methods return four values:

```go
(response *ResponseType, httpStatus int, success bool, errs []error)
```

Check `success` or `len(errs) > 0` to determine whether the call succeeded. The HTTP status code is always returned even on failure.

### STK Push (Lipa na M-Pesa Online)

Sends an STK push prompt to a customer's phone, requesting them to authorize a payment using their M-Pesa PIN.

```go
payload := &daraja.LipaNaMpesaOnlineRequestPayload{
    BusinessShortCode: "174379",
    Password:          "base64encodedpassword",
    Timestamp:         "20231205153000",
    TransactionType:   "CustomerPayBillOnline",
    Amount:            "100",
    PartyA:            "254708374149",
    PartyB:            "174379",
    PhoneNumber:       "254708374149",
    CallBackURL:       "https://example.com/mpesa/callback",
    AccountReference:  "ORDER-001",
    TransactionDesc:   "Payment for order",
}

response, status, ok, errs := client.LipaNaMpesaOnlinePayment(payload)
if ok {
    fmt.Println(response.MerchantRequestID)
    fmt.Println(response.CheckoutRequestID)
}
```

The `Password` field is a Base64 encoding of `BusinessShortCode + Passkey + Timestamp`.

### Business to Customer (B2C)

Sends money from a business shortcode to a customer's M-Pesa account. Use cases include salary disbursements, refunds, and promotions.

```go
payload := &daraja.B2CPaymentRequestPayload{
    InitiatorName:      "testapi",
    SecurityCredential: "encryptedcredential",
    CommandID:          "BusinessPayment",
    Amount:             "500",
    PartyA:             "600996",
    PartyB:             "254708374149",
    Remarks:            "Salary payment",
    QueueTimeOutURL:    "https://example.com/mpesa/timeout",
    ResultURL:          "https://example.com/mpesa/result",
    Occasion:           "Monthly salary",
}

response, status, ok, errs := client.B2CPaymentRequest(payload)
if ok {
    fmt.Println(response.ConversationID)
}
```

Accepted `CommandID` values: `SalaryPayment`, `BusinessPayment`, `PromotionPayment`.

### Business to Business (B2B)

Transfers funds between two business shortcodes.

```go
payload := &daraja.B2BPaymentRequestPayload{
    Initiator:              "testapi",
    SecurityCredential:     "encryptedcredential",
    CommandID:              "BusinessBuyGoods",
    SenderIdentifierType:   "4",
    RecieverIdentifierType: "4",
    Amount:                 "1000",
    PartyA:                 "600996",
    PartyB:                 "600000",
    AccountReference:       "INV-2023-001",
    Remarks:                "Payment for supplies",
    QueueTimeOutURL:        "https://example.com/mpesa/timeout",
    ResultURL:              "https://example.com/mpesa/result",
}

response, status, ok, errs := client.B2BPaymentRequest(payload)
if ok {
    fmt.Println(response.OriginatorConversationID)
}
```

### Register C2B URLs

Registers confirmation and validation URLs that Safaricom will call when a customer makes a C2B payment to your shortcode.

```go
payload := &daraja.RegisterURLRequestPayload{
    ShortCode:       "600996",
    ResponseType:    "Completed",
    ConfirmationURL: "https://example.com/mpesa/confirmation",
    ValidationURL:   "https://example.com/mpesa/validation",
}

response, status, ok, errs := client.RegisterURLs(payload)
if ok {
    fmt.Println(response.ResponseDescription)
}
```

`ResponseType` can be `Completed` (process all transactions even if validation times out) or `Cancelled` (cancel transactions that time out during validation).

### Query Transaction Status

Queries the current status of any M-Pesa transaction using its transaction ID.

```go
payload := &daraja.QueryTransactionStatusRequestPayload{
    Initiator:          "testapi",
    SecurityCredential: "encryptedcredential",
    CommandID:          "TransactionStatusQuery",
    TransactionID:      "OEI2AK4Q16",
    PartyA:             "600996",
    IdentifierType:     "4",
    ResultURL:          "https://example.com/mpesa/result",
    QueueTimeOutURL:    "https://example.com/mpesa/timeout",
    Remarks:            "Query",
    Occasion:           "",
}

response, status, ok, errs := client.QueryTransactionStatus(payload)
if ok {
    fmt.Println(response.ResponseDescription)
}
```

### Reverse a Transaction

Reverses a completed M-Pesa transaction. Only paybill and buy goods transactions can be reversed.

```go
payload := &daraja.ReversalRequestPayload{
    Initiator:              "testapi",
    SecurityCredential:     "encryptedcredential",
    CommandID:              "TransactionReversal",
    TransactionID:          "OEI2AK4Q16",
    Amount:                 "100",
    ReceiverParty:          "600996",
    RecieverIdentifierType: "4",
    ResultURL:              "https://example.com/mpesa/result",
    QueueTimeOutURL:        "https://example.com/mpesa/timeout",
    Remarks:                "Accidental payment",
    Occasion:               "",
}

response, status, ok, errs := client.ReverseTransaction(payload)
if ok {
    fmt.Println(response.ConversationID)
}
```

### Query Account Balance

Queries the available balance on a shortcode.

```go
payload := &daraja.AccountBalanceQueryRequestPayload{
    Initiatior:         "testapi",
    SecurityCredential: "encryptedcredential",
    CommandID:          "AccountBalance",
    PartyA:             "600996",
    IdentifierType:     "4",
    Remarks:            "Balance check",
    QueueTimeOutURL:    "https://example.com/mpesa/timeout",
    ResultURL:          "https://example.com/mpesa/result",
}

response, status, ok, errs := client.QueryAccountBalance(payload)
if ok {
    fmt.Println(response.ResponseDescription)
}
```

### Query Online Transaction Status

Queries the status of an STK push request using the `CheckoutRequestID` returned from `LipaNaMpesaOnlinePayment`.

The `OnlineTransactionQueryPayload` struct is provided for building the request payload. Use the same `Password` and `Timestamp` values as the original STK push request.

```go
payload := &daraja.OnlineTransactionQueryPayload{
    BusinessShortCode: "174379",
    Password:          "base64encodedpassword",
    Timestamp:         "20231205153000",
    CheckoutRequestID: "ws_CO_191220191020363925",
}
```

### Pull Transaction History

Pull APIs allow you to retrieve M-Pesa transaction history for a shortcode within a date range.

**Register for pull transactions:**

```go
registerPayload := &daraja.RegisterPullTransactionsRequestPayload{
    ShortCode:       "600996",
    RequestType:     "Pull",
    NominatedNumber: "254708374149",
    CallBackURL:     "https://example.com/mpesa/pull",
}
```

**Query pull transactions:**

```go
queryPayload := &daraja.QueryPullTransactionsRequestPayload{
    ShortCode:   "600996",
    StartDate:   "2023-12-01 00:00:00",
    EndDate:     "2023-12-31 23:59:59",
    OffSetValue: "0",
}
```

### C2B Validation

When Safaricom calls your validation URL before completing a C2B transaction, use `ValidateTransactionPayload` and its `ToResponse` method to accept or reject the payment.

```go
// In your HTTP handler for the validation URL:
func validationHandler(w http.ResponseWriter, r *http.Request) {
    var payload daraja.ValidateTransactionPayload
    json.NewDecoder(r.Body).Decode(&payload)

    // Validate the transaction (check account number, amount, etc.)
    valid := validateBusinessLogic(payload.BillRefNumber, payload.TransAmount)

    response := payload.ToResponse(daraja.ResultCodeInvalidAccount, valid)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

`ToResponse` sets `ResultDesc` to `"Accepted"` when `accept` is `true`, and `"Rejected"` when `false`. Use `ResultCode` `"0"` to accept, or one of the `ResultCode*` constants to reject with a specific reason.

Use `GetResultDesc` to get a human-readable description for any result code:

```go
desc := daraja.GetResultDesc(daraja.ResultCodeInvalidMSISDN)
// desc == "Invalid MSISDN"
```

---

## Error Handling

All API methods return a slice of errors as the last return value. A non-empty slice or a `false` success flag indicates failure.

```go
response, status, ok, errs := client.LipaNaMpesaOnlinePayment(payload)

if !ok {
    for _, err := range errs {
        log.Printf("error: %v", err)
    }
    log.Printf("http status: %d", status)
    return
}

// Use response
```

Errors in the slice may represent:

- Network failures
- JSON marshaling or unmarshaling failures
- HTTP request construction failures

The HTTP status code is returned regardless of whether the call succeeded. A status of `0` means the request could not be sent at all.

---

## Types Reference

### Core

| Type                  | Description                                                |
|-----------------------|------------------------------------------------------------|
| `Daraja`              | Client struct holding the access token and its expiry time |
| `DarajaAuthResponse`  | OAuth2 token response from the authorization endpoint      |

### STK Push

| Type                                       | Description                                                                        |
|--------------------------------------------|------------------------------------------------------------------------------------|
| `LipaNaMpesaOnlineRequestPayload`          | Request payload for STK push                                                       |
| `LipaNaMpesaOnlinePaymentResponsePayload`  | Response from STK push, includes `MerchantRequestID` and `CheckoutRequestID`       |
| `OnlineTransactionQueryPayload`            | Payload for querying the status of a pending STK push                              |

### B2C

| Type                        | Description                                              |
|-----------------------------|----------------------------------------------------------|
| `B2CPaymentRequestPayload`  | Request payload for business-to-customer payment         |
| `B2CPaymentResponsePayload` | Response including `ConversationID` and `ResponseCode`   |

### B2B

| Type                        | Description                                              |
|-----------------------------|----------------------------------------------------------|
| `B2BPaymentRequestPayload`  | Request payload for business-to-business payment         |
| `B2BPaymentResponsePayload` | Response including `ConversationID` and `ResponseCode`   |

### C2B

| Type                              | Description                                                           |
|-----------------------------------|-----------------------------------------------------------------------|
| `C2BConfirmationRequestPayload`   | Payload for C2B confirmation                                          |
| `RegisterURLRequestPayload`       | Payload for registering confirmation and validation URLs              |
| `RegisterURLResponsePayload`      | Response from URL registration                                        |
| `ValidateTransactionPayload`      | Incoming payload Safaricom sends to your validation URL               |
| `ValidationResponse`              | Your response to Safaricom, accepting or rejecting the transaction    |

### Transaction Management

| Type                                     | Description                                             |
|------------------------------------------|---------------------------------------------------------|
| `QueryTransactionStatusRequestPayload`   | Request payload for transaction status query            |
| `QueryTransactionStatusResponsePayload`  | Response from transaction status and reversal queries   |
| `ReversalRequestPayload`                 | Request payload for reversing a transaction             |
| `ReversalResponsePayload`                | Response from a reversal request                        |

### Account Balance

| Type                                | Description                                    |
|-------------------------------------|------------------------------------------------|
| `AccountBalanceQueryRequestPayload` | Request payload for account balance query      |
| `AccountBalanceQueryResponsePayload`| Response from account balance query            |

### Pull Transactions

| Type                                      | Description                                               |
|-------------------------------------------|-----------------------------------------------------------|
| `RegisterPullTransactionsRequestPayload`  | Request payload to register for pull transaction history  |
| `RegisterPullTransactionsResponsePayload` | Response from pull transaction registration               |
| `QueryPullTransactionsRequestPayload`     | Request payload to query transaction history              |
| `QueryPullTransactionsResponsePayload`    | Response from pull transaction query                      |

---

## Constants

Result code constants for C2B validation responses:

| Constant | Value | Description |
|----------|-------|-------------|
| `ResultCodeInvalidMSISDN` | `C2B00011` | Invalid MSISDN |
| `ResultCodeInvalidAccount` | `C2B00012` | Invalid Account Number |
| `ResultCodeInvalidAmount` | `C2B00013` | Invalid Amount |
| `ResultCodeInvalidKYC` | `C2B00014` | Invalid KYC Details |
| `ResultCodeInvalidShortcode` | `C2B00015` | Invalid Shortcode |
| `ResultCodeOtherError` | `C2B00016` | Other Error |

Use `GetResultDesc(code string) string` to look up the description for any of these codes at runtime.

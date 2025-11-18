## Daraja Go

```daraja-go``` is a minimal library to interact with the daraja API on safaricom

# Doc
```go
// Package daraja provides a lightweight Go client and data types for interacting
// with Safaricom's Daraja (MPesa) sandbox APIs. It includes helpers for
// authorization, making HTTP requests, and payload/response structs for common
// MPesa endpoints (STK push, B2B, B2C, account balance, reversals, transaction
// queries, and URL registration).

package daraja

// CONSUMER_SECRET holds the MPesa API consumer secret used to obtain an access
// token. It should be set by the caller (typically from configuration or
// environment).
var CONSUMER_SECRET = ""

// CONSUMER_KEY holds the MPesa API consumer key used to obtain an access
// token. It should be set by the caller (typically from configuration or
// environment).
var CONSUMER_KEY = ""

// SHORTCODE is a default business shortcode used by some MPesa operations.
// Populate it with the appropriate value for your account/environment.
var SHORTCODE = ""

// AUTH_URL is the endpoint used to obtain an OAuth access token from the
// Daraja sandbox.
const AUTH_URL = "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

// C2BConfirmation_URL is the (sandbox) endpoint for C2B confirmation/validation
// registration.
const C2BConfirmation_URL = "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"

// RegisterURL_URL is the (sandbox) endpoint for registering confirmation and
// validation URLs for C2B transactions.
const RegisterURL_URL = "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"

// AccountBalanceQuery_URL is the endpoint to query account balance.
const AccountBalanceQuery_URL = "https://sandbox.safaricom.co.ke/mpesa/accountbalance/v1/query"

// STK_PUSH_URL is the endpoint to initiate an STK push (Lipa Na Mpesa Online).
const STK_PUSH_URL = "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

// REVERSAL_URL is the endpoint to request a transaction reversal.
const REVERSAL_URL = "https://sandbox.safaricom.co.ke/mpesa/reversal/v1/request"

// B2B_URL is the endpoint to initiate B2B payment requests.
const B2B_URL = "https://sandbox.safaricom.co.ke/mpesa/b2b/v1/paymentrequest"

// TransactionStatusQuery_URL is the endpoint to query transaction status.
const TransactionStatusQuery_URL = "https://sandbox.safaricom.co.ke/mpesa/transactionstatus/v1/query"

// OnlineTransactionQuery_URL is the endpoint to query an STK push/online
// transaction status (checkout query).
const OnlineTransactionQuery_URL = "https://sandbox.safaricom.co.ke/mpesa/stkpushquery/v1/query"

// B2CPaymentRequest_URL is the endpoint to initiate B2C payment requests.
const B2CPaymentRequest_URL = "https://sandbox.safaricom.co.ke/mpesa/b2c/v1/paymentrequest"

// Daraja represents a minimal client holding the current access token and its
// expiry time. Methods on Daraja provide authorization and calls to various
// Daraja endpoints.
type Daraja struct {
    AccessToken string    // current OAuth access token
    Expiry      time.Time // UTC expiry time for the token
}

// DarajaAuthResponse models the OAuth token response returned by the Daraja
// authentication endpoint.
type DarajaAuthResponse struct {
    AccessToken string `json:"access_token"` // bearer token
    ExpiresIn   uint   `json:"expires_in"`   // lifetime in seconds
}

// post[T any] performs a POST request to the specified URL with the given JSON
// body and an optional bearer token. The generic type parameter T is the
// expected response payload type. Returns the HTTP status code, the unmarshalled
// payload of type T, and a slice of errors encountered (if any).
// Note: the function expects 'out' to be a struct value to determine type;
// this helper is internal and not exported.
func post[T any](url string, token string, body []byte, out T) (int, T, []error)

// get[T any] performs a GET request to the specified URL with an optional
// bearer token and attempts to unmarshal the JSON response into 'out'. Returns
// the HTTP status code, the unmarshalled payload of type T, and a slice of
// errors encountered (if any).
// Note: the function expects 'out' to be a struct value to determine type;
// this helper is internal and not exported.
func get[T any](url string, token string, out T) (int, T, []error)

// Authorize obtains a new access token from the Daraja AUTH_URL using the
// CONSUMER_KEY and CONSUMER_SECRET. On success it sets d.AccessToken and
// d.Expiry (UTC) and returns true with an empty error slice. On failure it
// returns false and a slice of errors detailing what went wrong.
func (d *Daraja) Authorize() (bool, []error)

// IsAuthorized reports whether the current Daraja access token is considered
// valid relative to the current UTC time. Intended to return true when the
// token is still valid and false when expired. Note: the implementation should
// be reviewed because an inverted comparison could report expired tokens as
// authorized.
func (d *Daraja) IsAuthorized() bool

// RetryAuth runs Authorize asynchronously in a goroutine and sends the result
// (authorized bool) to the provided status channel and any errors to the errs
// channel. Useful for non-blocking token refresh attempts.
func (d *Daraja) RetryAuth(status chan bool, errs chan []error)

// B2BPaymentRequest initiates a B2B payment using the provided payload. If the
// client is not authorized it will attempt to authorize first. Returns a
// pointer to the response payload, the HTTP status code, a boolean indicating
// whether the call was attempted (true on success path), and a slice of errors
// if any occurred.
func (d *Daraja) B2BPaymentRequest(r *B2BPaymentRequestPayload) (*B2BPaymentResponsePayload, int, bool, []error)

// ReverseTransaction requests a reversal for a given transaction. If the
// client is not authorized it will attempt to authorize first. Returns the
// parsed response payload pointer, HTTP status code, a boolean indicating
// whether the request was attempted, and any errors.
func (d *Daraja) ReverseTransaction(r *ReversalRequestPayload) (*QueryTransactionStatusResponsePayload, int, bool, []error)

// QueryTransactionStatus queries the status of a transaction. If the client is
// not authorized it will attempt to authorize first. Returns the parsed
// response payload pointer, HTTP status code, a boolean indicating whether the
// request was attempted, and any errors.
func (d *Daraja) QueryTransactionStatus(r *QueryTransactionStatusRequestPayload) (*QueryTransactionStatusResponsePayload, int, bool, []error)

// B2CPaymentRequest initiates a B2C payment using the provided payload. If the
// client is not authorized it will attempt to authorize first. Returns the
// parsed response payload pointer, HTTP status code, a boolean indicating
// whether the request was attempted, and any errors.
func (d *Daraja) B2CPaymentRequest(r *B2CPaymentRequestPayload) (*B2CPaymentResponsePayload, int, bool, []error)

// LipaNaMpesaOnlinePayment initiates an STK Push (Lipa Na Mpesa Online) request.
// If the client is not authorized it will attempt to authorize first. Returns
// the parsed response payload pointer, HTTP status code, a boolean indicating
// whether the request was attempted, and any errors.
func (d *Daraja) LipaNaMpesaOnlinePayment(r *LipaNaMpesaOnlineRequestPayload) (*LipaNaMpesaOnlinePaymentResponsePayload, int, bool, []error)

// QueryAccountBalance queries the account balance for the configured short
// code/party. If the client is not authorized it will attempt to authorize
// first. Returns the parsed response payload pointer, HTTP status code, a
// boolean indicating whether the request was attempted, and any errors.
func (d *Daraja) QueryAccountBalance(r *AccountBalanceQueryRequestPayload) (*AccountBalanceQueryResponsePayload, int, bool, []error)

// RegisterURLs registers confirmation and validation URLs for C2B transactions.
// If the client is not authorized it will attempt to authorize first. Returns
// the parsed response payload pointer, HTTP status code, a boolean indicating
// whether the request was attempted, and any errors.
func (d *Daraja) RegisterURLs(r *RegisterURLRequestPayload) (*RegisterURLResponsePayload, int, bool, []error)

// C2BConfirmationRequestPayload represents the JSON payload sent by Daraja
// (or expected for registration) for C2B confirmation and validation callbacks.
type C2BConfirmationRequestPayload struct {
    ShortCode     string `json:"ShortCode"`
    CommandID     string `json:"CommandID"`
    Amount        string `json:"Amount"`
    Msidsn        string `json:"Msisdn"`
    BillRefNumber string `json:"BillRefNumber"`
}

// B2BPaymentRequestPayload represents the request payload for a B2B payment.
type B2BPaymentRequestPayload struct {
    Initiator              string `json:"Initiator"`
    SecurityCredential     string `json:"SecurityCredential"`
    CommandID              string `json:"CommandID"`
    SenderIdentifierType   string `json:"SenderIdentifierType"`
    RecieverIdentifierType string `json:"RecieverIdentifierType"`
    Amount                 string `json:"Amount"`
    PartyA                 string `json:"PartyA"`
    PartyB                 string `json:"PartyB"`
    AccountReference       string `json:"AccountReference"`
    Remarks                string `json:"Remarks"`
    QueueTimeOutURL        string `json:"QueueTimeOutURL"`
    ResultURL              string `json:"ResultURL"`
}

// B2BPaymentResponsePayload models a typical response from a B2B payment
// request, including conversation IDs and a response code/description.
type B2BPaymentResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    ConversationID           string `json:"ConversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}

// RegisterURLRequestPayload is used to register confirmation and validation
// endpoints for C2B transactions.
type RegisterURLRequestPayload struct {
    ShortCode       string `json:"ShortCode"`
    ResponseType    string `json:"ResponseType"`
    ConfirmationURL string `json:"ConfirmationURL"`
    ValidationURL   string `json:"ValidationURL"`
}

// RegisterURLResponsePayload models the response from a Register URL request.
// Note: some responses from the API contain a misspelled field (OriginatorCoversationID)
// which is accounted for here.
type RegisterURLResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    OriginatorCoversationID  string `json:"OriginatorCoversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}

// AccountBalanceQueryRequestPayload represents the request payload for an
// account balance query operation.
type AccountBalanceQueryRequestPayload struct {
    Initiatior         string `json:"Initiator"`
    SecurityCredential string `json:"SecurityCredential"`
    CommandID          string `json:"CommandID"`
    PartyA             string `json:"PartyA"`
    IdentifierType     string `json:"IdentifierType"`
    Remarks            string `json:"Remarks"`
    QueueTimeOutURL    string `json:"QueueTimeOutURL"`
    ResultURL          string `json:"ResultURL"`
}

// AccountBalanceQueryResponsePayload models the response returned after an
// account balance query.
type AccountBalanceQueryResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    ConversationID           string `json:"ConversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}

// LipaNaMpesaOnlineRequestPayload represents the request payload for an STK
// push (Lipa Na Mpesa Online) operation.
type LipaNaMpesaOnlineRequestPayload struct {
    BusinessShortCode string `json:"BusinessShortCode"`
    Password          string `json:"Password"`
    Timestamp         string `json:"Timestamp"`
    TransactionType   string `json:"TransactionType"`
    Amount            string `json:"Amount"`
    PartyA            string `json:"PartyA"`
    PartyB            string `json:"PartyB"`
    PhoneNumber       string `json:"PhoneNumber"`
    CallBackURL       string `json:"CallBackURL"`
    AccountReference  string `json:"AccountReference"`
    TransactionDesc   string `json:"TransactionDesc"`
}

// LipaNaMpesaOnlinePaymentResponsePayload models the response returned after
// initiating an STK push request.
type LipaNaMpesaOnlinePaymentResponsePayload struct {
    MerchantRequestID   string `json:"MerchantRequestID"`
    CheckoutRequestID   string `json:"CheckoutRequestID"`
    ResponseCode        string `json:"ResponseCode"`
    ResponseDescription string `json:"ResponseDescription"`
    CustomerMessage     string `json:"CustomerMessage"`
}

// ReversalRequestPayload represents the payload to request a transaction
// reversal.
type ReversalRequestPayload struct {
    Initiator              string `json:"Initiator"`
    SecurityCredential     string `json:"SecurityCredential"`
    CommandID              string `json:"CommandID"`
    TransactionID          string `json:"TransactionID"`
    Amount                 string `json:"Amount"`
    ReceiverParty          string `json:"ReceiverParty"`
    RecieverIdentifierType string `json:"RecieverIdentifierType"`
    ResultURL              string `json:"ResultURL"`
    QueueTimeOutURL        string `json:"QueueTimeOutURL"`
    Remarks                string `json:"Remarks"`
    Occasion               string `json:"Occasion"`
}

// ReversalResponsePayload models the response returned after requesting a
// transaction reversal.
type ReversalResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    ConversationID           string `json:"ConversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}

// QueryTransactionStatusRequestPayload is used to query the status of a
// transaction.
type QueryTransactionStatusRequestPayload struct {
    Initiator          string `json:"Initiator"`
    SecurityCredential string `json:"SecurityCredential"`
    CommandID          string `json:"CommandID"`
    TransactionID      string `json:"TransactionID"`
    PartyA             string `json:"PartyA"`
    IdentifierType     string `json:"IdentifierType"`
    ResultURL          string `json:"ResultURL"`
    QueueTimeOutURL    string `json:"QueueTimeOutURL"`
    Remarks            string `json:"Remarks"`
    Occasion           string `json:"Occasion"`
}

// QueryTransactionStatusResponsePayload models the response to a transaction
// status query.
type QueryTransactionStatusResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    ConversationID           string `json:"ConversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}

// OnlineTransactionQueryPayload represents the payload required to query an
// STK push / online transaction by CheckoutRequestID.
type OnlineTransactionQueryPayload struct {
    BusinessShortCode string `json:"BusinessShortCode"`
    Password          string `json:"Password"`
    Timestamp         string `json:"Timestamp"`
    CheckoutRequestID string `json:"CheckoutRequestID"`
}

// B2CPaymentRequestPayload represents the payload for initiating a B2C
// payment (business-to-customer).
type B2CPaymentRequestPayload struct {
    InitiatorName      string `json:"InitiatorName"`
    SecurityCredential string `json:"SecurityCredential"`
    CommandID          string `json:"CommandID"`
    Amount             string `json:"Amount"`
    PartyA             string `json:"PartyA"`
    PartyB             string `json:"PartyB"`
    Remarks            string `json:"Remarks"`
    QueueTimeOutURL    string `json:"QueueTimeOutURL"`
    ResultURL          string `json:"ResultURL"`
    Occasion           string `json:"Occasion"`
}

// B2CPaymentResponsePayload models the response returned after initiating a
// B2C payment.
type B2CPaymentResponsePayload struct {
    OriginatorConversationID string `json:"OriginatorConversationID"`
    ConversationID           string `json:"ConversationID"`
    ResponseCode             string `json:"ResponseCode"`
    ResponseDescription      string `json:"ResponseDescription"`
}
```

Made with ❤️‍ by Sebastian Muchui
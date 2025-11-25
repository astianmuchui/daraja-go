package daraja

import (
	"time"

)

var (
	CONSUMER_SECRET = ""
	CONSUMER_KEY = ""
	SHORTCODE = ""
	PASSKEY = ""
	ACCOUNT_TYPE = ""
)

const (
	AUTH_URL                   = "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	C2BConfirmation_URL        = "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"
	RegisterURL_URL            = "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"
	AccountBalanceQuery_URL    = "https://sandbox.safaricom.co.ke/mpesa/accountbalance/v1/query"
	STK_PUSH_URL               = "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	REVERSAL_URL               = "https://sandbox.safaricom.co.ke/mpesa/reversal/v1/request"
	B2B_URL                    = "https://sandbox.safaricom.co.ke/mpesa/b2b/v1/paymentrequest"
	TransactionStatusQuery_URL = "https://sandbox.safaricom.co.ke/mpesa/transactionstatus/v1/query"
	OnlineTransactionQuery_URL = "https://sandbox.safaricom.co.ke/mpesa/stkpushquery/v1/query"
	B2CPaymentRequest_URL      = "https://sandbox.safaricom.co.ke/mpesa/b2c/v1/paymentrequest"
)

const (
	ResultCodeInvalidMSISDN    = "C2B00011"
	ResultCodeInvalidAccount   = "C2B00012"
	ResultCodeInvalidAmount    = "C2B00013"
	ResultCodeInvalidKYC       = "C2B00014"
	ResultCodeInvalidShortcode = "C2B00015"
	ResultCodeOtherError       = "C2B00016"
)

var ResultCodeDescriptions = map[string]string{
	ResultCodeInvalidMSISDN:    "Invalid MSISDN",
	ResultCodeInvalidAccount:   "Invalid Account Number",
	ResultCodeInvalidAmount:    "Invalid Amount",
	ResultCodeInvalidKYC:       "Invalid KYC Details",
	ResultCodeInvalidShortcode: "Invalid Shortcode",
	ResultCodeOtherError:       "Other Error",
}


type Daraja struct {
	AccessToken string
	Expiry      time.Time
}

type DarajaAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string   `json:"expires_in"`
}

type C2BConfirmationRequestPayload struct {
	ShortCode     string `json:"ShortCode"`
	CommandID     string `json:"CommandID"`
	Amount        string `json:"Amount"`
	Msidsn        string `json:"Msisdn"`
	BillRefNumber string `json:"BillRefNumber"`
}

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


type B2BPaymentResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type RegisterURLRequestPayload struct {
	ShortCode       string `json:"ShortCode"`
	ResponseType    string `json:"ResponseType"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ValidationURL   string `json:"ValidationURL"`
}

type RegisterURLResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	OriginatorCoversationID  string `json:"OriginatorCoversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
}

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

type AccountBalanceQueryResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

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

type LipaNaMpesaOnlinePaymentResponsePayload struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

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

type ReversalResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

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

type QueryTransactionStatusResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type OnlineTransactionQueryPayload struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

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

type B2CPaymentResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}



type ValidateTransactionPayload struct {
	TransactionType   string `json:"TransactionType"`
	TransID           string `json:"TransID"`
	TransTime         string `json:"TransTime"`
	TransAmount       string `json:"TransAmount"`
	BusinessShortCode string `json:"BusinessShortCode"`
	BillRefNumber     string `json:"BillRefNumber"`
	InvoiceNumber     string `json:"InvoiceNumber,omitempty"`
	OrgAccountBalance string `json:"OrgAccountBalance,omitempty"`
	ThirdPartyTransID string `json:"ThirdPartyTransID,omitempty"`
	MSISDN            string `json:"MSISDN"`
	FirstName         string `json:"FirstName,omitempty"`
	MiddleName        string `json:"MiddleName,omitempty"`
	LastName          string `json:"LastName,omitempty"`
}

type ValidationResponse struct {
	ResultCode string `json:"ResultCode"`
	ResultDesc string `json:"ResultDesc"`
}

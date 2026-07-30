/*
*
Daraja Go
By Sebastian Muchui :)
*/
package daraja

import (
	"time"
)

const (
	ResultCodeInvalidMSISDN    = "C2B00011"
	ResultCodeInvalidAccount   = "C2B00012"
	ResultCodeInvalidAmount    = "C2B00013"
	ResultCodeInvalidKYC       = "C2B00014"
	ResultCodeInvalidShortcode = "C2B00015"
	ResultCodeOtherError       = "C2B00016"
)

var prod bool

const (
	DARAJA_PRODUCTION_URL = "https://api.safaricom.co.ke"
	DARAJA_SANDBOX_URL    = "https://sandbox.safaricom.co.ke"
)

var (
	AUTH_URL                                 string
	C2BConfirmation_URL                      string
	RegisterURL_URL                          string
	AccountBalanceQuery_URL                  string
	STK_PUSH_URL                             string
	REVERSAL_URL                             string
	B2B_URL                                  string
	TransactionStatusQuery_URL               string
	OnlineTransactionQuery_URL               string
	B2CPaymentRequest_URL                    string
	TransactionHistoryQuery_URL              string
	TransactionHistoryRegister_URL           string
	TaxRemittanceRequest_URL                 string
	BusinessPaybillTransaction_URL           string
	BusinessBuyGoodsTransaction_URL          string
	BillManagerGenericOptIn_URL              string
	BillManagerSingleInvoicingGeneric_URL    string
	BillManagerBulkInvoicingGeneric_URL      string
	BillManagerPaymentsAndReconciliation_URL string
	BillManagerCancelSingleInvoice_URL       string
	BillManagerCancelBulkInvoices_URL        string
	BillManagerUpdateOptinDetailsRequest_URL string

	B2BExpressCheckout_URL string
	BusinessToPochi_URL    string
)

func Production(state bool) {
	prod = state
	initializeURLs()
}

func SetProductionMode() {
	Production(true)
}

func initializeURLs() {
	var url_prefix string
	if !prod {
		url_prefix = "https://sandbox.safaricom.co.ke"
	} else {
		url_prefix = "https://api.safaricom.co.ke"
	}

	AUTH_URL = url_prefix + "/oauth/v1/generate?grant_type=client_credentials"
	C2BConfirmation_URL = url_prefix + "/mpesa/c2b/v1/registerurl"

	RegisterURL_URL = url_prefix + "/mpesa/c2b/v1/registerurl"
	AccountBalanceQuery_URL = url_prefix + "/mpesa/accountbalance/v1/query"

	STK_PUSH_URL = url_prefix + "/mpesa/stkpush/v1/processrequest"
	REVERSAL_URL = url_prefix + "/mpesa/reversal/v1/request"

	B2B_URL = url_prefix + "/mpesa/b2b/v1/paymentrequest"
	TransactionStatusQuery_URL = url_prefix + "/mpesa/transactionstatus/v1/query"

	OnlineTransactionQuery_URL = url_prefix + "/mpesa/stkpushquery/v1/query"
	B2CPaymentRequest_URL = url_prefix + "/mpesa/b2c/v1/paymentrequest"

	TransactionHistoryRegister_URL = url_prefix + "/pulltransactions/v1/register"
	TransactionHistoryQuery_URL = url_prefix + "/pulltransactions/v1/query"

	TaxRemittanceRequest_URL = url_prefix + "/mpesa/b2b/v1/remittax"

	BusinessPaybillTransaction_URL = url_prefix + "/mpesa/b2b/v1/paymentrequest"
	BusinessBuyGoodsTransaction_URL = BusinessPaybillTransaction_URL

	BillManagerGenericOptIn_URL = url_prefix + "v1/billmanager-invoice/optin"
	BillManagerSingleInvoicingGeneric_URL = url_prefix + "v1/billmanager-invoice/single-invoicing"

	BillManagerBulkInvoicingGeneric_URL = url_prefix + "v1/billmanager-invoice/bulk-invoicing"
	BillManagerPaymentsAndReconciliation_URL = url_prefix + "v1/billmanager-invoice/reconciliation"

	BillManagerCancelSingleInvoice_URL = url_prefix + "/v1/billmanager-invoice/cancel-single-invoice"
	BillManagerCancelBulkInvoices_URL = url_prefix + "/v1/billmanager-invoice/cancel-bulk-invoices"

	BillManagerUpdateOptinDetailsRequest_URL = url_prefix + "/v1/billmanager-invoice/change-optin-details "

	B2BExpressCheckout_URL = url_prefix + "/v1/ussdpush/get-msisdn"
	BusinessToPochi_URL = url_prefix + "/mpesa/b2pochi/v1/paymentrequest"
}

func init() {
	prod = false
	initializeURLs()
}

var (
	CONSUMER_SECRET = ""
	CONSUMER_KEY    = ""
	SHORTCODE       = ""
	PASSKEY         = ""
	ACCOUNT_TYPE    = ""
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
	ExpiresIn   string `json:"expires_in"`
}

type C2BConfirmationRequestPayload struct {
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

type B2BPaymentRequestPayload struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"Command ID"`
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

type RegisterPullTransactionsRequestPayload struct {
	ShortCode       string `json:"ShortCode"`
	RequestType     string `json:"RequestType"`
	NominatedNumber string `json:"NominatedNumber"`
	CallBackURL     string `json:"CallBackURL"`
}

type RegisterPullTransactionsResponsePayload struct {
	ResponseRefID       string `json:"ResponseRefID"`
	ResponseStatus      string `json:"ResponseStatus"`
	ShortCode           string `json:"ShortCode"`
	ResponseDescription string `json:"ResponseDescription"`
}

type QueryPullTransactionsRequestPayload struct {
	ShortCode   string `json:"ShortCode"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
	OffSetValue string `json:"OffSetValue"`
}

type QueryPullTransactionsResponsePayload struct {
	ResponseStatus      string `json:"ResponseStatus"`
	ResponseDescription string `json:"ResponseDescription"`
	ResponseRefID       string `json:"ResponseRefID"`
	ShortCode           string `json:"ShortCode"`
	StartDate           string `json:"StartDate"`
	EndDate             string `json:"EndDate"`
	OffSetValue         string `json:"OffSetValue"`
}

type RemitKRARequestPayload struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"Command ID"`
	SenderIdentifierType   string `json:"SenderIdentifierType"`
	ReceiverIdentifierType string `json:"RecieverIdentifierType"`
	Amount                 string `json:"Amount"`
	PartyA                 string `json:"PartyA"`
	PartyB                 string `json:"PartyB"`
	AccountReference       string `json:"AccountReference"`
	Remarks                string `json:"Remarks"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	ResultURL              string `json:"ResultURL"`
}

type RemitKRAResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type RemitTaxResultPayload struct {
	Result struct {
		ResultType               string `json:"ResultType"`
		ResultCode               string `json:"ResultCode"`
		ResultDesc               string `json:"ResultDesc"`
		OriginatorConversationID string `json:"OriginatorConversationID"`
		ConversationID           string `json:"ConversationID"`
		TransactionID            string `json:"TransactionID"`
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ResultParameter"`
		} `json:"ResultParameters"`
		ReferenceData struct {
			ReferenceItem []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ReferenceItem"`
		} `json:"ReferenceData"`
	} `json:"Result"`
}

type RemitTaxFailedResultPayload struct {
	Result struct {
		ResultType               string `json:"ResultType"`
		ResultCode               int    `json:"ResultCode"`
		ResultDesc               string `json:"ResultDesc"`
		OriginatorConversationID string `json:"OriginatorConversationID"`
		ConversationID           string `json:"ConversationID"`
		TransactionID            string `json:"TransactionID"`
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ResultParameter"`
		} `json:"ResultParameters"`
		ReferenceData struct {
			ReferenceItem struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ReferenceItem"`
		} `json:"ReferenceData"`
	} `json:"Result"`
}

type BusinessToPaybillTransactionRequestPayload struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"Command ID"`
	SenderIdentifierType   string `json:"SenderIdentifierType"`
	RecieverIdentifierType string `json:"RecieverIdentifierType"`
	Amount                 string `json:"Amount"`
	PartyA                 string `json:"PartyA"`
	PartyB                 string `json:"PartyB"`
	AccountReference       string `json:"AccountReference"`
	Requester              string `json:"Requester"`
	Remarks                string `json:"Remarks"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	ResultURL              string `json:"ResultURL"`
}

type BusinessToBuyGoodsTransactionRequestPayload struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"Command ID"`
	SenderIdentifierType   string `json:"SenderIdentifierType"`
	RecieverIdentifierType string `json:"RecieverIdentifierType"`
	Amount                 string `json:"Amount"`
	PartyA                 string `json:"PartyA"`
	PartyB                 string `json:"PartyB"`
	AccountReference       string `json:"AccountReference"`
	Requester              string `json:"Requester"`
	Remarks                string `json:"Remarks"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	ResultURL              string `json:"ResultURL"`
}

type GenericResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type GenericResult struct {
	Result struct {
		ResultType               string `json:"ResultType"`
		ResultCode               string `json:"ResultCode"`
		ResultDesc               string `json:"ResultDesc"`
		OriginatorConversationID string `json:"OriginatorConversationID"`
		ConversationID           string `json:"ConversationID"`
		TransactionID            string `json:"TransactionID"`
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ResultParameter"`
		} `json:"ResultParameters"`
		ReferenceData struct {
			ReferenceItem []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ReferenceItem"`
		} `json:"ReferenceData"`
	} `json:"Result"`
}

type BusinessToPaybillTransactionResponsePayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type BillManagerOptInRequestPayload struct {
	ShortCode       string `json:"shortcode"`
	Email           string `json:"email"`
	OfficialContact string `json:"officialContact"`
	SendReminders   string `json:"sendReminders"`
	Logo            []byte `json:"logo"`
	CallbackURL     string `json:"callbackurl"`
}

type BillManagerOptInResponsePayload struct {
	AppKey  string `json:"app_key"`
	ResMsg  string `json:"resmsg"`
	ResCode string `json:"rescode"`
}

type BillManagerSingleInvoiceRequestPayload struct {
	ExternalReference string        `json:"externalReference"`
	BilledFullName    string        `json:"billedFullName"`
	BilledPhoneNumber string        `json:"billedPhoneNumber"`
	BilledPeriod      string        `json:"billedPeriod"`
	InvoiceName       string        `json:"invoiceName"`
	DueDate           string        `json:"dueDate"`
	AccountReference  string        `json:"accountReference"`
	Amount            string        `json:"amount"`
	InvoiceItems      []InvoiceItem `json:"invoiceItems"`
}

type InvoiceItem struct {
	ItemName string `json:"itemName"`
	Amount   string `json:"amount"`
}

type BillManagerResponsePayload struct {
	StatusMessage string `json:"Status_Message"`
	ResMsg        string `json:"resmsg"`
	ResCode       string `json:"rescode"`
}

type BulkInvoicingGenericApiRequestPayload struct {
	Invoices []BillManagerSingleInvoiceRequestPayload `json:"invoices"`
}

type BillManagerPaymentReconcilRequestPayload struct {
	TransactionID    string `json:"transactionId"`
	PaidAmount       string `json:"paidAmount"`
	MSISDN           string `json:"msisdn"`
	DateCreated      string `json:"dateCreated"`
	AccountReference string `json:"accountReference"`
	ShortCode        string `json:"shortCode"`
}

type BillManagerReceiptAcknowledgementPayload struct {
	PaymentDate       string `json:"paymentDate"`
	PaidAmount        string `json:"paidAmount"`
	AccountReference  string `json:"accountReference"`
	TransactionID     string `json:"transactionId"`
	PhoneNumber       string `json:"phoneNumber"`
	FullName          string `json:"fullName"`
	InvoiceName       string `json:"invoiceName"`
	ExternalReference string `json:"externalReference"`
}
type BillManagerPaymentAcknowledgmentResultPayload struct {
	StatusCode       string `json:"statusCode"`
	StatusMessage    string `json:"statusMessage"`
	TransactionID    string `json:"transactionId"`
	AccountReference string `json:"accountReference"`
}

type BillManagerPaymentsReconGenericResponse struct {
	ResMsg  string `json:"resmsg"`
	ResCode string `json:"rescode"`
}

type BillManagerCancelInvoiceRequestPayload struct {
	ExternalReference string `json:"externalReference"`
}

type BillManagerCancelInvoiceResponsePayload struct {
	StatusMessage string        `json:"Status_Message"`
	ResMsg        string        `json:"resmsg"`
	ResCode       string        `json:"rescode"`
	Errors        []interface{} `json:"errors"`
}

type BillManagerUpdateOptInDetailsRequestPayload struct {
	BillManagerOptInRequestPayload
}

type B2BEpressCheckoutRequestPayload struct {
	PrimaryShortCode  string `json:"primaryShortCode"`
	ReceiverShortCode string `json:"receiverShortCode"`
	Amount            string `json:"amount"`
	PaymentRef        string `json:"paymentRef"`
	CallbackUrl       string `json:"callbackUrl"`
	PartnerName       string `json:"partnerName"`
	RequestRefID      string `json:"RequestRefID"`
}

type B2BExpressCheckoutResponsePayload struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type USSDCallbackResponsePayload struct {
	ResultCode       string `json:"resultCode"`
	ResultDesc       string `json:"resultDesc"`
	RequestID        string `json:"requestId"`
	Amount           string `json:"amount"`
	PaymentReference string `json:"paymentReference"`
	ResultType       string `json:"resultType,omitempty"`
	ConversationID   string `json:"conversationID,omitempty"`
	TransactionID    string `json:"transactionId,omitempty"`
	Status           string `json:"status,omitempty"`
}


type B2PochiPaymentRequestPayload struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	InitiatorName           string `json:"InitiatorName"`
	SecurityCredential      string `json:"SecurityCredential"`
	CommandID               string `json:"CommandID"`
	Amount                  string `json:"Amount"`
	PartyA                  string `json:"PartyA"`
	PartyB                  string `json:"PartyB"`
	Remarks                 string `json:"Remarks"`
	QueueTimeOutURL         string `json:"QueueTimeOutURL"`
	ResultURL               string `json:"ResultURL"`
	Occasion                string `json:"Occasion"`
}

type CallbackPayload struct {
	Result struct {
		ResultType               string `json:"ResultType"`
		ResultCode               string `json:"ResultCode"`
		ResultDesc               string `json:"ResultDesc"`
		OriginatorConversationID string `json:"OriginatorConversationID"`
		ConversationID           string `json:"ConversationID"`
		TransactionID            string `json:"TransactionID"`
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ResultParameter"`
		} `json:"ResultParameters"`
		ReferenceData struct {
			ReferenceItem []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"ReferenceItem"`
		} `json:"ReferenceData"`
	} `json:"Result"`
}


type B2PochiErrorResponse struct {
	RequestID    string `json:"requestId"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}


type B2CAccountTopUpRequestPayload struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"CommandID"`
	SenderIdentifierType   string `json:"SenderIdentifierType"`
	RecieverIdentifierType string `json:"RecieverIdentifierType"`
	Amount                 string `json:"Amount"`
	PartyA                 string `json:"PartyA"`
	PartyB                 string `json:"PartyB"`
	AccountReference       string `json:"AccountReference"`
	Requester              string `json:"Requester"`
	Remarks                string `json:"Remarks"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	ResultURL              string `json:"ResultURL"`
}
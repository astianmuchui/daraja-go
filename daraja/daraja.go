package daraja

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func post[T any](url string, token string, body []byte, out T) (int, T, []error) {

	if reflect.TypeOf(out).Kind() != reflect.Struct {
		return 0, out, []error{}
	}

	var errs []error
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return 0, out, []error{err}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, out, []error{err}
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, out, []error{err}
	}

	status := resp.StatusCode

	var payload T

	err = json.Unmarshal(res, &payload)

	if err != nil {
		errs = append(errs, err)
	}

	/*
	   Daraja reports a rejected request with an HTTP status and a body like

	       {"requestId":"...","errorCode":"...","errorMessage":"..."}

	   which unmarshals cleanly into any of the response types here, because
	   none of them declare those fields. The result was a zero-valued
	   response, no error, and the reason discarded: a caller saw an empty
	   struct and had nothing to log. Surface the body so the reason survives.
	*/
	if status < 200 || status > 299 {
		errs = append(errs, fmt.Errorf("daraja: %s responded %d: %s", url, status, bodySnippet(res)))
	}

	return status, payload, errs
}

// bodySnippet trims a response body to something loggable.
func bodySnippet(body []byte) string {
	text := strings.TrimSpace(string(body))
	if text == "" {
		return "(empty body)"
	}
	if len(text) > 400 {
		return text[:400]
	}
	return text
}

func SendRequest[X any](d *Daraja, url string, token string, body []byte, out X) (int, X, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	if reflect.TypeOf(out).Kind() != reflect.Struct {
		return 0, *new(X), []error{}
	}
	status, payload, errs := post(url, token, body, out)

	return status, payload, errs
}

func get[T any](url string, token string, out T) (int, T, []error) {
	if reflect.TypeOf(out).Kind() != reflect.Struct {
		return 0, out, []error{}
	}

	var errs []error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, out, []error{err}
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, out, []error{err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, out, []error{err}
	}

	err = json.Unmarshal(body, out)

	if err != nil {
		errs = append(errs, err)
	}

	return resp.StatusCode, out, errs
}

func (d *Daraja) Authorize() (bool, []error) {
	var errs []error
	var status int
	var body []byte

	req, err := http.NewRequest("GET", AUTH_URL, nil)
	req.SetBasicAuth(CONSUMER_KEY, CONSUMER_SECRET)

	if err != nil {
		errs = append(errs, err)
	} else {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err)
		} else {
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				errs = append(errs, err)
			}
			status = resp.StatusCode
		}
	}

	resTime := time.Now().UTC()

	if len(errs) > 0 {
		return false, errs
	}

	if status != 200 {
		return false, errs
	}

	var resp DarajaAuthResponse

	err = json.Unmarshal(body, &resp)

	if err != nil {
		log.Println("Error Unmarshalling Body", err)
		errs = append(errs, err)
		return false, errs
	}

	expiry, err := strconv.Atoi(resp.ExpiresIn)

	if err != nil {
		log.Println("Error Unmarshalling Body", err)
		errs = append(errs, err)
		return false, errs
	}

	d.AccessToken = resp.AccessToken
	d.Expiry = resTime.Add(time.Duration(expiry) * time.Second)
	d.Expiry = d.Expiry.UTC()

	return true, []error{}
}

/*
	deternines whether or not the auth token is
	valid by checking the expiry against the current time
*/

func (d *Daraja) IsAuthorized() bool {
	return time.Now().UTC().After(d.Expiry.UTC())
}

func (d *Daraja) RetryAuth(status chan bool, errs chan []error) {
	go func() {
		authorized, errors := d.Authorize()
		status <- authorized
		errs <- errors
	}()
}

func (d *Daraja) B2BPaymentRequest(r *B2BPaymentRequestPayload) (*B2BPaymentResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &B2BPaymentResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(B2B_URL, d.AccessToken, payload, B2BPaymentResponsePayload{})

	if len(errs) > 0 {
		return &B2BPaymentResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) ReverseTransaction(r *ReversalRequestPayload) (*ReversalResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &ReversalResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(REVERSAL_URL, d.AccessToken, payload, ReversalResponsePayload{})

	if len(errs) > 0 {
		return &ReversalResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) QueryTransactionStatus(r *QueryTransactionStatusRequestPayload) (*QueryTransactionStatusResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &QueryTransactionStatusResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(TransactionStatusQuery_URL, d.AccessToken, payload, QueryTransactionStatusResponsePayload{})

	if len(errs) > 0 {
		return &QueryTransactionStatusResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) B2CPaymentRequest(r *B2CPaymentRequestPayload) (*B2CPaymentResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &B2CPaymentResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(B2CPaymentRequest_URL, d.AccessToken, payload, B2CPaymentResponsePayload{})

	if len(errs) > 0 {
		return &B2CPaymentResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) LipaNaMpesaOnlinePayment(r *LipaNaMpesaOnlineRequestPayload) (*LipaNaMpesaOnlinePaymentResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &LipaNaMpesaOnlinePaymentResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(STK_PUSH_URL, d.AccessToken, payload, LipaNaMpesaOnlinePaymentResponsePayload{})

	if len(errs) > 0 {
		return &LipaNaMpesaOnlinePaymentResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) QueryAccountBalance(r *AccountBalanceQueryRequestPayload) (*AccountBalanceQueryResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &AccountBalanceQueryResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(AccountBalanceQuery_URL, d.AccessToken, payload, AccountBalanceQueryResponsePayload{})

	if len(errs) > 0 {
		return &AccountBalanceQueryResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) RegisterURLs(r *RegisterURLRequestPayload) (*RegisterURLResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &RegisterURLResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(RegisterURL_URL, d.AccessToken, payload, RegisterURLResponsePayload{})

	if len(errs) > 0 {
		return &RegisterURLResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func GetResultDesc(code string) string {
	return ResultCodeDescriptions[code]
}

func (v *ValidateTransactionPayload) ToResponse(ResultCode string, accept bool) *ValidationResponse {
	var r ValidationResponse
	r.ResultCode = ResultCode

	var status string

	if accept {
		status = "Accepted"
	} else {
		status = "Rejected"
	}

	r.ResultDesc = status

	return &r
}

func (d *Daraja) RemitTax(r *RemitKRARequestPayload) (*RemitKRAResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &RemitKRAResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(TaxRemittanceRequest_URL, d.AccessToken, payload, RemitKRAResponsePayload{})

	if len(errs) > 0 {
		return &RemitKRAResponsePayload{}, status, false, errs
	}

	return &response, status, true, []error{}
}

func (d *Daraja) BusinessToPaybillTransaction(r *BusinessToPaybillTransactionRequestPayload) (*GenericResponse, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &GenericResponse{}, 0, false, []error{err}
	}

	status, response, errs := post(BusinessPaybillTransaction_URL, d.AccessToken, payload, GenericResponse{})
	if len(errs) > 0 {
		return &GenericResponse{}, status, false, errs
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BusinessToBuyGoodsTransaction(r *BusinessToBuyGoodsTransactionRequestPayload) (*GenericResponse, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return &GenericResponse{}, 0, false, []error{err}
	}

	status, response, errors := post(BusinessBuyGoodsTransaction_URL, d.AccessToken, payload, GenericResponse{})

	if len(errors) > 0 {
		return &GenericResponse{}, status, false, errors
	}

	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerOptIn(r *BillManagerOptInRequestPayload) (*BillManagerOptInResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerOptInResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerGenericOptIn_URL, d.AccessToken, payload, BillManagerOptInResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerOptInResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerSendSingleInvoice(r *BillManagerSingleInvoiceRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerSingleInvoicingGeneric_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerSendBulkInvoice(r *[]BillManagerSingleInvoiceRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerBulkInvoicingGeneric_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerReconcilePayment(r *BillManagerPaymentReconcilRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerPaymentsAndReconciliation_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

/*
Nobody:
Payment: "Acknowledge Meeee :)"
*/

func (d *Daraja) BillManagerAcknowledgePayment(r *BillManagerReceiptAcknowledgementPayload, BillManagerReceiptAcknowledgement_URL string) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerReceiptAcknowledgement_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerCancelSingleInvoice(r *BillManagerCancelInvoiceRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerCancelSingleInvoice_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerBulkCancelInvoices(r *[]BillManagerCancelInvoiceRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerCancelBulkInvoices_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) BillManagerUpdateOptInDetails(r *BillManagerUpdateOptInDetailsRequestPayload) (*BillManagerResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &BillManagerResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(BillManagerUpdateOptinDetailsRequest_URL, d.AccessToken, payload, BillManagerResponsePayload{})
	if len(errors) > 0 {
		return &BillManagerResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) InitiateB2BExpressCheckout(r *B2BEpressCheckoutRequestPayload) (*B2BExpressCheckoutResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return &B2BExpressCheckoutResponsePayload{}, 0, false, []error{err}
	}

	status, response, errors := post(B2BExpressCheckout_URL, d.AccessToken, payload, B2BExpressCheckoutResponsePayload{})
	if len(errors) > 0 {
		return &B2BExpressCheckoutResponsePayload{}, status, false, errors
	}
	return &response, status, true, []error{}
}

func (d *Daraja) InitiateB2PochiPayment(r *B2PochiPaymentRequestPayload) (*GenericResponse, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return &GenericResponse{}, 0, false, []error{err}
	}

	status, response, errors := post(BusinessToPochi_URL, d.AccessToken, payload, GenericResponse{})
	if len(errors) > 0 {
		return &GenericResponse{}, status, false, errors
	}

	return &response, status, true, []error{}
}

func (d *Daraja) InitiateB2CAccountTopUp(r *B2CAccountTopUpRequestPayload) (*GenericResponse, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)
	if err != nil {
		return &GenericResponse{}, 0, false, []error{err}
	}

	status, response, errors := post(BusinessToPochi_URL, d.AccessToken, payload, GenericResponse{})

	if len(errors) > 0 {
		return &GenericResponse{}, status, false, errors
	}

	return &response, status, true, []error{}
}

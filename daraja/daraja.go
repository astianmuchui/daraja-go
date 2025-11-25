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
	"time"

)

func post [T any] (url string, token string, body []byte, out T) (int, T, []error) {

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

	return status, payload, errs
}

func get [T any] (url string, token string, out T) (int, T, []error) {
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
	return time.Now().UTC().After(d.Expiry)
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

func (d *Daraja) ReverseTransaction(r *ReversalRequestPayload) (*QueryTransactionStatusResponsePayload, int, bool, []error) {
	if !d.IsAuthorized() {
		d.Authorize()
	}

	payload, err := json.Marshal(r)

	if err != nil {
		return &QueryTransactionStatusResponsePayload{}, 0, false, []error{err}
	}

	status, response, errs := post(REVERSAL_URL, d.AccessToken, payload, QueryTransactionStatusResponsePayload{})

	if len(errs) > 0 {
		return &QueryTransactionStatusResponsePayload{}, status, false, errs
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

func (d *Daraja) LipaNaMpesaOnlinePayment(r *LipaNaMpesaOnlineRequestPayload)  (*LipaNaMpesaOnlinePaymentResponsePayload, int, bool, []error) {
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

func (v *ValidateTransactionPayload) ToResponse(ResultCode string, accept bool) (*ValidationResponse) {
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

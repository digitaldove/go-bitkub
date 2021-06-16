package bitkub

import "fmt"

// https://github.com/bitkub/bitkub-official-api-docs/blob/master/restful-api.md#error-codes
var errorMessages = make(map[int]string)

func init() {
	errorMessages[0] = "No error" // should not be an error
	errorMessages[1] = "Invalid JSON payload"
	errorMessages[2] = "Missing X-BTK-APIKEY"
	errorMessages[3] = "Invalid API key"
	errorMessages[4] = "API pending for activation"
	errorMessages[5] = "IP not allowed"
	errorMessages[6] = "Missing / invalid signature"
	errorMessages[7] = "Missing timestamp"
	errorMessages[8] = "Invalid timestamp"
	errorMessages[9] = "Invalid user"
	errorMessages[10] = "Invalid parameter"
	errorMessages[11] = "Invalid symbol"
	errorMessages[12] = "Invalid amount"
	errorMessages[13] = "Invalid rate"
	errorMessages[14] = "Improper rate"
	errorMessages[15] = "Amount too low"
	errorMessages[16] = "Failed to get balance"
	errorMessages[17] = "Wallet is empty"
	errorMessages[18] = "Insufficient balance"
	errorMessages[19] = "Failed to insert order into db"
	errorMessages[20] = "Failed to deduct balance"
	errorMessages[21] = "Invalid order for cancellation"
	errorMessages[22] = "Invalid side"
	errorMessages[23] = "Failed to update order status"
	errorMessages[24] = "Invalid order for lookup"
	errorMessages[25] = "KYC level 1 is required to proceed"
	errorMessages[30] = "Limit exceeds"
	errorMessages[40] = "Pending withdrawal exists"
	errorMessages[41] = "Invalid currency for withdrawal"
	errorMessages[42] = "Address is not in whitelist"
	errorMessages[43] = "Failed to deduct crypto"
	errorMessages[44] = "Failed to create withdrawal record"
	errorMessages[45] = "Nonce has to be numeric"
	errorMessages[46] = "Invalid nonce"
	errorMessages[47] = "Withdrawal limit exceeds"
	errorMessages[48] = "Invalid bank account"
	errorMessages[49] = "Bank limit exceeds"
	errorMessages[50] = "Pending withdrawal exists"
	errorMessages[51] = "Withdrawal is under maintenance"
	errorMessages[90] = "Server error (please contact support)"
}

type btkError struct {
	Code    int
	Message string
}

func (err btkError) Error() string {
	msg, ok := errorMessages[err.Code]
	if !ok {
		msg = fmt.Sprintf("Unknown error %d", err.Code)
	}
	if err.Message != "" {
		return fmt.Sprintf("%s: %s", msg, err.Message)
	}
	return msg
}

func newBtkError(code int, message ...string) (err btkError) {
	err.Code = code
	if len(message) > 0 {
		err.Message = message[0]
	}
	return err
}

package SMS

type SMSBase interface {
	SendSMS(sms SMS) (*SMS, error)
}

type SMS struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

package SMS

type SMSBase interface{
	SendSMS(sms SMS) (*SMS, error)
}

type SMS struct{
	Phone string
	Code string
}

const numberBytes = "0123456789"


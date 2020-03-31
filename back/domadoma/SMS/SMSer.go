package SMS

import (
	"../../domadoma"
	"math/rand"
	"net/http"
)

func NewSMSBase(configfile *domadoma.ConfigFile) (SMSBase, error) {
	req, err := http.NewRequest("GET","https://smsc.kz/sys/send.php",nil)
	q := req.URL.Query()
	q.Add("login", configfile.SMSlogin)
	q.Add("psw", configfile.SMSpass)
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}
	return &requestSMSBase{req: req}, nil
}

type requestSMSBase struct {
	req *http.Request
}

func (f *requestSMSBase)SendSMS(sms SMS) (*SMS, error){
	client := &http.Client{}
	newreq := f.req
	q := f.req.URL.Query()
	sms.Code = RandStringBytes(6)
	q.Add("phones","7"+sms.Phone)
	q.Add("mes","edimdoma kod:"+sms.Code)
	newreq.URL.RawQuery = q.Encode()

	_, err := client.Do(newreq)
	if err != nil {
		return nil,err
	}
	return &sms,nil
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Int63()%int64(len(numberBytes))]
	}
	return string(b)
}


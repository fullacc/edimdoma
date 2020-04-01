package SMS

import (
	crand "crypto/rand"
	"encoding/binary"
	"github.com/fullacc/edimdoma/back/domadoma"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	sms.Code = RandCode()
	q.Add("phones","7"+sms.Phone)
	q.Add("mes","edimdoma kod:"+sms.Code)
	newreq.URL.RawQuery = q.Encode()

	_, err := client.Do(newreq)
	if err != nil {
		return nil,err
	}
	return &sms,nil
}

func RandCode() string {
	var src cryptoSource
	rnd := rand.New(src)
	return strconv.Itoa(100000+rnd.Intn(899999))
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
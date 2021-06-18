package model

type EventBody struct {
	CallId       string `json:"call_id"`
	SipCallId    string `json:"sip_call_id"`
	Domain       string `json:"domain"`
	Direction    string `json:"direction"`
	FromNumber   string `json:"from_number"`
	ToNumber     string `json:"to_number"`
	Hotline      string `json:"hotline"`
	State        string `json:"state"`
	Duration     int    `json:"duration"`
	Billsec      int    `json:"billsec"`
	RecordingUrl string `json:"recording_url"`
}

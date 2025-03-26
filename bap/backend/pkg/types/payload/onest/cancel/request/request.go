package request

type SeekerCancelPayload struct {
	WorkerID string              `json:"worker_id,omitempty"`
	BppID   string              `json:"bpp_id,omitempty"`
	BppURI  string              `json:"bpp_uri,omitempty"`
	Location SelectedJobLocation `json:"location,omitempty"`
}

type SelectedJobLocation struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	AreaCode    string      `json:"area_code"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type CancelRequest struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}
type City struct {
	Code string `json:"code"`
}
type Country struct {
	Code string `json:"code"`
}
type Location struct {
	City    City    `json:"city"`
	Country Country `json:"country"`
}
type Context struct {
	Domain        string   `json:"domain"`
	Action        string   `json:"action"`
	Version       string   `json:"version"`
	BapID         string   `json:"bap_id"`
	BapURI        string   `json:"bap_uri"`
	BppID         string   `json:"bpp_id"`
	BppURI        string   `json:"bpp_uri"`
	TransactionID string   `json:"transaction_id"`
	MessageID     string   `json:"message_id"`
	Location      Location `json:"location"`
	Timestamp     string   `json:"timestamp"`
	TTL           string   `json:"ttl"`
}
type Message struct {
	OrderID              string `json:"order_id"`
	CancellationReasonID string `json:"cancellation_reason_id"`
}

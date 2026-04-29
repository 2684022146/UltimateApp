package model

type GeocodeResponse struct {
	Status   string    `json:"status"`
	Geocodes []Geocode `json:"geocodes"`
}
type Geocode struct {
	Location string `json:"location"`
}

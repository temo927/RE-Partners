package transport

type PackSizesResponse struct {
	Sizes []int `json:"sizes"`
}

type UpdatePackSizesRequest struct {
	Sizes []int `json:"sizes"`
}

type CalculateRequest struct {
	Items int `json:"items"`
}

type PackResponse struct {
	Size     int `json:"size"`
	Quantity int `json:"quantity"`
}

type CalculateResponse struct {
	Packs []PackResponse `json:"packs"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

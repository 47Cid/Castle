package podAPI

type Response struct {
	StatusCode int
	Valid      string `json:"isValid"`
}

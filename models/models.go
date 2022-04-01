package models

type Survivor struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Gender   string   `json:"gender"`
	Location Location `json:"location"`
	Resource []string `json:"resource"`
	Status   int      `json:"status"` //infected status >= 3
}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

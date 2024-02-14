package model

type MainCategory struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Thumbnail   string         `json:"thumbnail"`
	SubCategory *[]MainCategory `json:"sub_category,omitempty"`
}

	
type Message struct {
	Code    int    `json:"code"  validate:"required"`
	Message string `json:"message"  validate:"required"`
}

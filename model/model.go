package model

type AssetList struct {
	ID            string         `json:"id" validate:"required,email" `
	Name          string         `json:"name" validate:"required" `
	Thumbnail     string         `json:"thumbnail" validate:"required" `
	SubCategories []SubAssetList `json:"sub_categories" validate:"required"`
}

type SubAssetList struct {
	ID        string `json:"id" validate:"required,email" `
	Name      string `json:"name" validate:"required" `
	Thumbnail string `json:"thumbnail" validate:"required" `
}

type Message struct {
	Code    int    `json:"code"  validate:"required"`
	Message string `json:"message"  validate:"required"`
}

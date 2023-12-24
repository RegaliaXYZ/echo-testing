package models

type CreateModelInput struct {
	Name       string `json:"model" binding:"required"`
	Service    string `json:"service" binding:"required"`
	SubService string `json:"sub_service"`
	Language   string `json:"language" binding:"required"`
}

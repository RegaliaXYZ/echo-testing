package models

// Create Model
type CreateModelInput struct {
	Name       string `json:"model" binding:"required"`
	Service    string `json:"service" binding:"required"`
	SubService string `json:"sub_service"`
	Language   string `json:"language" binding:"required"`
}

// Upload Content

type GlobalUploadContentInput struct {
	Payload  []Payload `json:"payload" binding:"required"`
	Complete string    `json:"complete" binding:"required"`
}

type Payload struct {
	Utterance string `json:"utt" binding:"required"`
	Intent    string `json:"intent,omitempty"`
	Num       int    `json:"num,omitempty"`
	Tags      []struct {
		Name        string `json:"name" binding:"required"`
		StartOffset int    `json:"start" binding:"required"`
		EndOffset   int    `json:"end" binding:"required"`
	} `json:"tags,omitempty"`
	Type string `json:"type,omitempty"`
}

type NLUPayload struct {
	Utterance string `json:"utt" binding:"required"`
	Intent    string `json:"intent" binding:"required"`
	Num       int    `json:"num" binding:"required"`
	Type      string `json:"type,omitempty"`
}

type UploadContentNLUInput struct {
	Payload  []NLUPayload `json:"payload" binding:"required"`
	Complete string       `json:"complete" binding:"required"`
}

type NERPayload struct {
	Utterance string `json:"utt" binding:"required"`
	Tags      []struct {
		Name        string `json:"name" binding:"required"`
		StartOffset int    `json:"start" binding:"required"`
		EndOffset   int    `json:"end" binding:"required"`
	} `json:"tags" binding:"required"`
	Type string `json:"type,omitempty"`
}

type UploadContentNERInput struct {
	Payload  []NERPayload `json:"payload" binding:"required"`
	Complete string       `json:"complete" binding:"required"`
}

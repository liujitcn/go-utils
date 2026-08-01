package volc

type translateTextRequest struct {
	SourceLanguage string   `json:"SourceLanguage,omitempty"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}

type translation struct {
	Text                   string `json:"Translation"`
	DetectedSourceLanguage string `json:"DetectedSourceLanguage,omitempty"`
}

type responseError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type responseMetadata struct {
	RequestID string         `json:"RequestId"`
	Action    string         `json:"Action"`
	Version   string         `json:"Version"`
	Service   string         `json:"Service"`
	Region    string         `json:"Region"`
	Error     *responseError `json:"Error,omitempty"`
}

type translateTextResponse struct {
	ResponseMetadata responseMetadata `json:"ResponseMetadata"`
	ResponseMetaData responseMetadata `json:"ResponseMetaData"`
	TranslationList  []translation    `json:"TranslationList"`
}

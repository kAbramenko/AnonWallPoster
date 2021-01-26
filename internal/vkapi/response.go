package vkapi

type VkError struct {
	ErrorCode     float64       `json:"error_code"`
	ErrorMsg      string        `json:"error_msg"`
	RequestParams []interface{} `json:"request_params"`
}
type VkResp struct {
	AlbumID   float64 `json:"album_id"`
	UploadURL string  `json:"upload_url"`
	UserID    float64 `json:"user_id"`
	Server    float64 `json:"server"`
	Photo     string  `json:"photo"`
	Hash      string  `json:"hash"`
}
type VkResponse struct {
	// Response map[string]interface{} `json:"response"`
	Response VkResp `json:"response"`
	// Error    map[string]interface{} `json:"error"`
	Error VkError `json:"error"`
}

type VkResponseStats struct {
	Response []struct {
		Activity struct {
			Comments     int
			Copies       int
			Hidden       int
			Likes        int
			Subscribed   int
			UnSubscribed int
		}
		Visitors struct {
			Views      int
			Visitors   int
			MobileView int `json:"mobile_views"`
		}
	}
}
type VkResponseArray struct {
	Response []VkPhotoObject `json:"response"`
}

type VkPhotoObject struct {
	AlbumID   int                      `json:"album_id"`
	Date      int                      `json:"date"`
	ID        int                      `json:"id"`
	OwnerID   int                      `json:"owner_id"`
	HasTags   bool                     `json:"has_tags"`
	AccessKey string                   `json:"access_key"`
	Sizes     []map[string]interface{} `json:"sizes"`
	Text      string                   `json:"text"`
}

type VkUploadResponse struct {
	Server int    `json:"server"`
	Hash   string `json:"hash"`
	Photo  string `json:"photo"`
}

// HasError return true if error filed has been filled
func (v *VkResponse) HasError() bool {
	if v.Error.ErrorMsg == "" {
		return false
	}

	return true
}

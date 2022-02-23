package model

type (
	PostResponse struct {
		Data PostResponseData
	}

	PostResponseData struct {
		Children []PostResponseDataChildren
	}

	PostResponseDataChildren struct {
		Data PostResponseDataChildrenData
	}

	PostResponseDataChildrenData struct {
		Title   string `json:"link_title"`
		LinkUrl string `json:"link_url"`
	}
)

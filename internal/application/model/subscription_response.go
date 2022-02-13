package model

type (
	SubscriptionsResponse struct {
		Data SubscriptionsDateResponse `json:"data"`
	}

	SubscriptionsDateResponse struct {
		After    string                          `json:"after"`
		Children []SubscriptionsChildrenResponse `json:"children"`
	}

	SubscriptionsChildrenResponse struct {
		Kind string                            `json:"kind"`
		Data SubscriptionsChildrenDateResponse `json:"data"`
	}

	SubscriptionsChildrenDateResponse struct {
		DispalyName         string `json:"display_name"`
		DispalyNamePrefixed string `json:"display_name_prefixed"`
	}
)

package apimodels

type OrderCreationRequests struct {
	UserID      string        `json:"user_id" binding:"required"`
	AuthToken   string        `json:"auth_token" binding:"required"`
	SellerID    string        `json:"seller_id"`
	Attachments []string      `json:"attachments" binding:"required"`
	Services    []ServicePair `json:"services" binding:"required"`
	Comment     *string       `json:"comment" binding:"required"`
}

type ServicePair struct {
	First  string `json:"first" binding:"required"`
	Second int    `json:"second" binding:"required"`
}

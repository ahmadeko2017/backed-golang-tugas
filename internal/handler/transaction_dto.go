package handler

type CheckoutItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type CheckoutRequest struct {
	Items []CheckoutItemRequest `json:"items" binding:"required,dive,required"`
	Total float64               `json:"total" binding:"required"`
}

type ReportResponse struct {
	TotalRevenue      float64     `json:"total_revenue"`
	TotalTransactions int64       `json:"total_transactions"`
	BestSeller        *TopProduct `json:"best_seller,omitempty"`
}

type TopProduct struct {
	Name    string `json:"name"`
	SoldQty int64  `json:"sold_qty"`
}

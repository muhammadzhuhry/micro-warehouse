package response

import "micro-warehouse/product-service/pkg/pagination"

type ProductResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Barcode    string           `json:"barcode"`
	About      string           `json:"about"`
	Price      float64          `json:"price"`
	IsPopular  bool             `json:"is_popular"`
	CategoryID uint             `json:"category_id"`
	Thumbnail  string           `json:"thumbnail"`
	Category   CategoryResponse `json:"category"`
}

type GetAllProductsResponse struct {
	Products   []ProductResponse             `json:"products"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

package request

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required"`
	Barcode    string `json:"barcode" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Thumbnail  string `json:"thumbnail" validate:"required"`
	About      string `json:"about" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	IsPopular  bool   `json:"is_popular" validate:"required"`
}

type GetAllProductsRequest struct {
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search    string `query:"search" validate:"omitempty"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=id name created_at"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

package response

import "micro-warehouse/product-service/pkg/pagination"

type CategoryResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Tagline       string `json:"tagline"`
	Photo         string `json:"photo"`
	CountProducts int    `json:"count_products"`
}

type GetAllCategoriesResponse struct {
	Categories []CategoryResponse            `json:"categories"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

package request

type CreateCategoryRequest struct {
	Name    string `json:"name" validate:"required"`
	Tagline string `json:"tagline" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllCategoriesRequest struct {
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search    string `query:"search" validate:"omitempty"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=id name created_at"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

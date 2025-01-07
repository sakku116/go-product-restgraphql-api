// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package product_gql

type CreateProductReq struct {
	Name  string `json:"name"`
	Stock string `json:"stock"`
	Price string `json:"price"`
}

type CreateProductRespData struct {
	Data *Product `json:"data"`
}

type DeleteProductRespData struct {
	Data *Product `json:"data"`
}

type GetProductByUUIDRespData struct {
	Data *Product `json:"data"`
}

type GetProductListReq struct {
	UserUUID  *string `json:"userUUID,omitempty"`
	Query     *string `json:"query,omitempty"`
	QueryBy   *string `json:"queryBy,omitempty"`
	Page      *int    `json:"page,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
	SortOrder *int    `json:"sortOrder,omitempty"`
	SortBy    *string `json:"sortBy,omitempty"`
}

type GetProductListRespData struct {
	Data        []*Product `json:"data"`
	Total       string     `json:"total"`
	CurrentPage string     `json:"currentPage"`
	TotalPage   string     `json:"totalPage"`
}

type Mutation struct {
}

type Product struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"userUUID"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Stock     string `json:"stock"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Query struct {
}

type UpdateProductReq struct {
	Name  *string `json:"name,omitempty"`
	Stock *string `json:"stock,omitempty"`
	Price *string `json:"price,omitempty"`
}

type UpdateProductRespData struct {
	Data *Product `json:"data"`
}

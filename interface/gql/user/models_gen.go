// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package user_gql

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateUserRespData struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Email     string `json:"email"`
}

type DeleteUserRespData struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetUserByUUIDResp struct {
	UUID      string `json:"uuid"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}

type GetUserListReq struct {
	Query     *string `json:"query,omitempty"`
	QueryBy   *string `json:"queryBy,omitempty"`
	Page      *string `json:"page,omitempty"`
	Limit     *string `json:"limit,omitempty"`
	SortOrder *string `json:"sortOrder,omitempty"`
	SortBy    *string `json:"sortBy,omitempty"`
}

type GetUserListRespData struct {
	Total       string                     `json:"total"`
	CurrentPage string                     `json:"currentPage"`
	TotalPage   string                     `json:"totalPage"`
	Data        []*GetUserListRespDataUser `json:"data"`
}

type GetUserListRespDataUser struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateUserReq struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
}

type UpdateUserRespData struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

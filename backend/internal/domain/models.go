package domain

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type List struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserPermissions struct {
	UserID      uint   `json:"user_id"`
	ListID      uint   `json:"list_id"`
	Permissions uint16 `json:"permissions"`
}

type Level struct {
	ID           uint   `json:"id"`
	ListID       *uint  `json:"list_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ShowcaseLink string `json:"showcase_link"`
	Placement    int    `json:"placement"`
}

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

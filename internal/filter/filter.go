package filter

type Filter struct {
	Order     string `json:"order,omitempty"`
	Username  string `json:"username,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
	Page      int    `json:"page,omitempty"`
	Page_Size int    `json:"page_size,omitempty"`
}

func (f Filter) CurrentPage() int {
	return (f.Page - 1) * f.Page_Size
}

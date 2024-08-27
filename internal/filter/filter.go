package filter

import "math"

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

type MetaData struct {
	Page         int
	PageSize     int
	FirstPage    int
	LastPage     int
	TotalRecords int
}

func CalculateMetaData(totalRecords, page, pageSize int) MetaData {
	if totalRecords == 0 {
		return MetaData{}
	}

	return MetaData{
		Page:         page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

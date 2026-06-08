package pagination

import "github.com/gin-gonic/gin"

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Result struct {
	Items    interface{} `json:"items"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
}

func FromQuery(c *gin.Context) Page {
	page := parseInt(c.DefaultQuery("page", "1"), 1)
	size := parseInt(c.DefaultQuery("page_size", "20"), 20)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 20
	}
	return Page{Page: page, PageSize: size}
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func parseInt(s string, fallback int) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return fallback
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

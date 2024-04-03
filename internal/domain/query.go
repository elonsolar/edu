package domain

/*
*

	{
		"expr":{


		},
		"page":{
			"pageNo":1,
			"pageSize":10
		}


	}

*
*/

type Expression struct {
	Op      string // and,or , = like
	SubExpr []*Expression

	// normal
	IsLogic bool
	Column  string
	Value   interface{} // "1" ,1 , [1,2,3], ["1","2","3"]
}

// omitempty struct ->json 忽略零值
type Page struct {
	PageNo   int `json:"page_no,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

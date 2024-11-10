package model

//查询详情, 所有详情都通过id和唯一键进行详情查询
type DESCRIBE_BY int64

const (
    DESCRIBE_BY_ID   DESCRIBE_BY = 0 // 通过id查询
    DESCRIBE_BY_NAME DESCRIBE_BY = 1 // 通过name查询
)

// 查询详情请求 1.根据id查询 2.根据name查询
type DetailOptions struct {
    DescribeBy DESCRIBE_BY `json:"describe_by"`
    Id         string      `json:"id,omitempty"`
    Name       string      `json:"name,omitempty"`
}

func NewDescribeRequestByName(name string) *DetailOptions {
    return &DetailOptions{
        DescribeBy: DESCRIBE_BY_NAME,
        Name:       name,
    }
}

func NewDescribeRequestByID(id string) *DetailOptions {
    return &DetailOptions{
        DescribeBy: DESCRIBE_BY_ID,
        Id:         id,
    }
}

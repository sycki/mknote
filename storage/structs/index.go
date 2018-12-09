package structs

//该结构体用在首页
//name表示一个tag
//articles表示属于该tag的文章名
type ArticleTag struct {
	Name     string
	Articles []*Article
}

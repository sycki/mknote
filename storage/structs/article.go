package structs

type ArticleFull struct {
	id           int
	parent_id    int
	title        string
	en_name      string
	content      string
	author       int
	create_date  string
	change_date  string
	status       string
	tags         string
	like_count   int
	unlike_count int
	viewer_count int
}

type Article struct {
	Id           string
	Title        string
	En_name      string
	Content      string
	Author       string
	Like_count   int
	Viewer_count int
	Create_date  string
}

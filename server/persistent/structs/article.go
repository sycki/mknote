/*
Copyright 2017 sycki.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

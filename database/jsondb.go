package database

import (
	"encoding/json"
	"html/template"
	"log"
	"strings"
	"sycki/structs"

	"github.com/go-redis/redis"
)

var (
	DBcli     *redis.Client
	Templates map[string]*template.Template
)

func init() {
	DBcli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Templates = make(map[string]*template.Template)
	log.Println("INFO init jsonDAO complete.")
}

func JGET(key string, path string) (string, error) {
	return DBcli.JGET(key, path).Result()
}

// field format: id__titel__tag__publish
func Index() (resutl string, err error) {
	keys, _, err := DBcli.Scan(0, "*__publish", 50).Result()
	if err != nil {
		return
	}
	arr := []*structs.ArticleTag{}
	m := make(map[string][]string)
	for _, key := range keys {
		tag := strings.Split(key, "__")[2]
		m[tag] = append(m[tag], tag)
	}
	for k, v := range m {
		arr = append(arr, &structs.ArticleTag{k, v})
	}
	result, err := json.Marshal(arr)
	return string(result), err
}

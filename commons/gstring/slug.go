package gstring

import "github.com/gosimple/slug"

func Slug(str string) string {
	return slug.Make(str)
}

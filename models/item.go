package models

import "encoding/json"

type Icon struct {
	Path string `json:"path"`
}

type Item struct {
	Arg      string `json:"arg"`
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Icon     Icon   `json:"icon"`
}

func NewItem(title, subtitle, arg string, path string) *Item {
	return &Item{
		Arg:      arg,
		Title:    title,
		SubTitle: subtitle,
		Icon:     Icon{Path: path},
	}
}

type Items struct {
	Items []*Item `json:"items"`
}

func NewItems() *Items {
	return &Items{Items: make([]*Item, 0)}
}

func (i *Items) Append(item *Item) *Items {
	i.Items = append(i.Items, item)
	return i
}

func (i *Items) Encode() string {
	dat, err := json.Marshal(*i)
	if err != nil {
		return ""
	}

	return string(dat)
}

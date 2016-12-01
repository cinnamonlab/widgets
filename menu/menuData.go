package menu

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"sort"
	"reflect"
)

type Items struct {
	Items []Item `json:"items"`
}

var storeItems *Items

type Item struct {
	Id        int64 `json:"id"`
	ParentId  int64 `json:"parent_id"`
	Name      string    `json:"name"`
	Url       string `json:"url"`
	HtmlClass string `json:"html_class"`
	Icon      string `json:"icon"`
}

func StoreItems() *Items {
	if storeItems == nil {
		storeItems = &Items{}
	}
	return storeItems
}

func GetItems1() []Item {
	var items []Item
	store := StoreItems()
	if len(store.Items) > 0 {

	} else {
		pwd, _ := os.Getwd()
		raw, err := ioutil.ReadFile(pwd + "/widgets/menu/items.json")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		json.Unmarshal(raw, &items)
		store.Items = items
	}
	return store.Items
}
func GetItems() []Item {
	var items []Item
	//raw, err := ioutil.ReadFile("./widgets/menu/items.json")
	raw, err := ioutil.ReadFile("./adminmenus.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &items)
	By(Prop("ParentId", true)).Sort(items)
	return items
}

type itemSorter struct {
	items []Item
	by    func(p1, p2 *Item) bool // Closure used in the Less method.
}
// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j])
}

type By func(p1, p2 *Item) bool

func (by By) Sort(items []Item) {
	item := &itemSorter{
		items: items,
		by: by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(item)
}

func Prop(field string, asc bool) func(p1, p2 *Item) bool {
	return func(p1, p2 *Item) bool {

		v1 := reflect.Indirect(reflect.ValueOf(p1)).FieldByName(field)
		v2 := reflect.Indirect(reflect.ValueOf(p2)).FieldByName(field)

		ret := false

		switch v1.Kind() {
		case reflect.Int64:
			ret = int64(v1.Int()) < int64(v2.Int())
		case reflect.Float64:
			ret = float64(v1.Float()) < float64(v2.Float())
		case reflect.String:
			ret = string(v1.String()) < string(v2.String())
		}

		if asc {
			return ret
		}
		return !ret
	}
}
package menu

import (
	"bytes"
	"html/template"
	"github.com/cinnamonlab/gorest"
)

var  menuActive = "active"
type Menu struct {
	Items []*menuItem
}

type menuItem struct {
	Id int64
	ParentId int64
	Name string
	Url string
	HtmlClass string
	Icon string
	Items     []*menuItem
}

func New() *Menu {
	return new(Menu)
}

func (self *Menu) Init() *Menu {
	items := GetItems()
	self.Items = self.Items[:0]
	if items != nil {
		for _, item := range items {
			if item.ParentId == 0 {
				self.AddItem(item)
			}
			if item.ParentId > 0 {
				self.AddSubItem(item.ParentId, item)
			}
		}
	}
	return self
}

func (self *Menu) AddItem(item Item) *Menu {
	currentUrl := gorest.GetCurrentPath()
	itemMenu := &menuItem{item.Id, item.ParentId, item.Name, item.Url, item.HtmlClass, item.Icon, nil}
	if currentUrl == item.Url {
		itemMenu.HtmlClass = menuActive
	}
	self.Items = append(self.Items, itemMenu)
	return self
}

func (self *Menu) AddSubItem(parentId int64, subItem Item) *Menu {
	currentUrl := gorest.GetCurrentPath()
	for key, item := range self.Items {
		if item.Id != parentId {
			continue
		}
		itemMenu := &menuItem{subItem.Id, subItem.ParentId, subItem.Name, subItem.Url, subItem.HtmlClass, subItem.Icon, nil}
		if currentUrl == subItem.Url {
			self.Items[key].HtmlClass = menuActive
			itemMenu.HtmlClass = menuActive
		}
		self.Items[key].Items = append(self.Items[key].Items, itemMenu)
	}
	return self
}

func (self *Menu) Sort() {

}

func (self *Menu) Render() template.HTML {
	self.Sort()
	// Create a new template and parse the letter into it.
	var out bytes.Buffer
	tMenu := template.Must(template.New("menu").Parse(tmplMenu))
	tMap := map[string]interface{}{
		"menu": self,
	}
	tMenu.Execute(&out, tMap)
	return template.HTML(out.String())
}

const tmplMenu = `
<ul class="sidebar-menu">
  {{ range $index, $element := .menu.Items }}
    {{ if len $element.Items }}
      <li class="treeview {{$element.HtmlClass}}">
          <a href="{{$element.Url}}"><i class="fa fa-{{$element.Icon}}"></i><span>{{$element.Name}}</span> <i class="fa fa-angle-left pull-right"></i></a>
          <ul class="treeview-menu menu-open">
            {{ range $element.Items }}
              <li class="{{.HtmlClass}}"><a href="{{.Url}}"><i class="fa fa-{{.Icon}}"></i> {{.Name}}</a></li>
            {{ end }}
          </ul>
      </li>
    {{ else }}
      <li class="treeview {{$element.HtmlClass}}"><a href="{{$element.Url}}"><i class="fa fa-{{$element.Icon}}"></i><span>{{$element.Name}}</span></a></li>
    {{ end }}
  {{ end }}
</ul>
`

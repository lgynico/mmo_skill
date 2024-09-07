package conf

import (
	"fmt"
	"path"
	"reflect"
)

/* adapter */
type Conf interface {
	getFileName() string
	load(fullPath string)
}

type ConfAdapter interface {
	onLoadComplete()
}

type BaseConfAdapter struct {
}

func (a *BaseConfAdapter) onLoadComplete() {

}

/* adapter set */
type adapterSet struct {
	*HeroConfAdapter
	*SkillConfAdapter
}

/* confMgr */
type confMgr struct {
	*confHolder
	path          string
	adpaterKeySet map[string]bool
	// mapName2Conf    map[string]Conf
	// mapName2Adapter map[string]ConfAdapter
	// adapters        []ConfAdapter
}

var ConfMgr = &confMgr{
	path:          ".",
	adpaterKeySet: make(map[string]bool),
	confHolder: &confHolder{
		mapName2Conf: make(map[string]Conf),
		adapters:     make([]ConfAdapter, 0),
		adapterSet:   &adapterSet{},
		confSet:      &confSet{},
	},
	// mapName2Conf:    make(map[string]Conf),
	// mapName2Adapter: make(map[string]ConfAdapter),
	// adapters:        make([]ConfAdapter, 0),
}

func (c *confMgr) Init(confPath string) {
	c.path = confPath
	c.loadConfigs(c.path)
}

func (c *confMgr) Reload() {
	confHolder := &confHolder{
		mapName2Conf: make(map[string]Conf, len(c.mapName2Conf)),
		adapters:     make([]ConfAdapter, 0, len(c.adapters)),
		adapterSet:   &adapterSet{},
		confSet:      &confSet{},
		// mapName2Adapter: make(map[string]ConfAdapter, len(c.mapName2Adapter)),
	}

	confHolder.copy(c.confHolder)
	confHolder.loadConfigs(c.path)
	c.confHolder = confHolder
}

func (c *confMgr) addConfMap(k string, v Conf) {
	if _, ok := c.mapName2Conf[k]; ok {
		panic(fmt.Sprintf("config error: config %s had existed", k))
	}

	c.mapName2Conf[k] = v
	s := reflect.TypeOf(v).Elem().Name()
	field := reflect.ValueOf(c).Elem().FieldByName(s)
	field.Set(reflect.ValueOf(v))
}

func (c *confMgr) addAdapter(k string, adapter ConfAdapter) {
	if _, ok := c.adpaterKeySet[k]; ok {
		panic(fmt.Sprintf("config error: adapter %s exists.", k))
	}
	c.adpaterKeySet[k] = true
	c.adapters = append(c.adapters, adapter)

	s := reflect.TypeOf(adapter).Elem().Name()
	field := reflect.ValueOf(c).Elem().FieldByName(s)
	field.Set(reflect.ValueOf(adapter))
}

type confHolder struct {
	*adapterSet
	*confSet
	mapName2Conf map[string]Conf
	adapters     []ConfAdapter
	// mapName2Adapter map[string]ConfAdapter
}

func (c *confHolder) loadConfigs(filePath string) {
	for _, cfg := range c.mapName2Conf {
		fileName := cfg.getFileName()
		fullPath := path.Join(filePath, fileName)
		cfg.load(fullPath)
	}

	for _, adapter := range c.adapters {
		adapter.onLoadComplete()
	}
}

func (c *confHolder) copy(src *confHolder) {
	for k, conf := range src.mapName2Conf {
		t := reflect.TypeOf(conf).Elem()
		v := reflect.New(t)
		c.mapName2Conf[k] = v.Interface().(Conf)

		s := t.Name()
		field := reflect.ValueOf(c).Elem().FieldByName(s)
		field.Set(v)
	}

	for _, adapter := range src.adapters {
		t := reflect.TypeOf(adapter).Elem()
		v := reflect.New(t)
		c.adapters = append(c.adapters, v.Interface().(ConfAdapter))

		s := t.Name()
		field := reflect.ValueOf(c).Elem().FieldByName(s)
		field.Set(v)
	}
}

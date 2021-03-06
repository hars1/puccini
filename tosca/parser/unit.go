package parser

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca"
)

//
// NoEntity
//

type NoEntity struct {
	Context *tosca.Context
}

func NewNoEntity(toscaContext *tosca.Context) *NoEntity {
	return &NoEntity{toscaContext}
}

// tosca.Contextual interface
func (self *NoEntity) GetContext() *tosca.Context {
	return self.Context
}

//
// Unit
//

type Unit struct {
	EntityPtr       interface{}
	Container       *Unit
	Imports         Units
	NameTransformer tosca.NameTransformer
	Locker          sync.Mutex
}

func NewUnitNoEntity(toscaContext *tosca.Context, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	return NewUnit(NewNoEntity(toscaContext), container, nameTransformer)
}

func NewUnit(entityPtr interface{}, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	self := Unit{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
	}
	if container != nil {
		container.AddImport(&self)
	}
	return &self
}

func (self *Unit) AddImport(import_ *Unit) {
	self.Locker.Lock()
	self.Imports = append(self.Imports, import_)
	self.Locker.Unlock()
}

func (self *Unit) GetContext() *tosca.Context {
	return tosca.GetContext(self.EntityPtr)
}

// Print

func (self *Unit) PrintImports(indent int, treePrefix terminal.TreePrefix) {
	length := len(self.Imports)
	last := length - 1

	// Sort
	imports := make(Units, length)
	copy(imports, self.Imports)
	sort.Sort(imports)

	for i, unit := range imports {
		isLast := i == last
		unit.PrintNode(indent, treePrefix, isLast)
		unit.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *Unit) PrintNode(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorValue(self.GetContext().URL.String()))
}

//
// Units
//

type Units []*Unit

// sort.Interface

func (self Units) Len() int {
	return len(self)
}

func (self Units) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Units) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}

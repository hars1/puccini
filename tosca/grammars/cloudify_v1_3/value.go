package cloudify_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Description *string
	Type        string

	rendered bool
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	ToFunctionCalls(context)
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderParameter(dataType, nil, false)
	} else {
		self.RenderParameter(dataType, definition.ParameterDefinition, false)
	}
}

func (self *Value) RenderParameter(dataType *DataType, definition *ParameterDefinition, allowNil bool) {
	if self.rendered {
		// Avoid rendering more than once (can happen if we use the "default" value)
		return
	}
	self.rendered = true

	if (definition != nil) && (definition.Description != nil) {
		self.Description = definition.Description
	} else {
		self.Description = dataType.Description
	}

	self.Type = dataType.Name

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// TODO: dataType.Complete(self.Context.Data, self.Context)

	// Internal types
	if internalTypeName, ok := dataType.GetInternalTypeName(); ok {
		if validator, ok := tosca.PrimitiveTypeValidators[internalTypeName]; ok {
			if self.Context.Data == nil {
				// Nil data only happens when an parameter is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Context.Data = tosca.PrimitiveTypeZeroes[internalTypeName]
			}

			// Primitive types
			if !validator(self.Context.Data) {
				self.Context.ReportValueWrongType(dataType.Name)
			}
		} else {
			// Special types
			if read, ok := self.Context.Grammar.Readers[internalTypeName]; ok {
				self.Context.Data = read(self.Context)
			} else {
				// Avoid reporting more than once
				if !dataType.typeProblemReported {
					dataType.Context.ReportUnsupportedType()
					dataType.typeProblemReported = true
				}
			}
		}
	} else if self.Context.ValidateType("map") {
		// Complex data types

		map_ := self.Context.Data.(ard.Map)

		// All properties must be defined in type
		for key := range map_ {
			name := yamlkeys.KeyString(key)
			if _, ok := dataType.PropertyDefinitions[name]; !ok {
				self.Context.MapChild(name, nil).ReportUndeclared("property")
				delete(map_, key)
			}
		}

		// Render properties
		for key, definition := range dataType.PropertyDefinitions {
			if data, ok := map_[key]; ok {
				var value *Value
				if value, ok = data.(*Value); !ok {
					// Convert to value
					value = NewValue(self.Context.MapChild(key, data))
					map_[key] = value
				}
				if definition.DataType != nil {
					value.RenderProperty(definition.DataType, definition)
				}
			} else {
				// PropertyDefinition.Required defaults to true
				required := (definition.Required == nil) || *definition.Required
				if required {
					self.Context.MapChild(key, data).ReportPropertyRequired("property")
				}
			}
		}
	}
}

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.List:
		list := normal.NewList(len(data))
		for index, value := range data {
			list.Set(index, NewValue(self.Context.ListChild(index, value)).Normalize())
		}
		constrainable = list

	case ard.Map:
		map_ := normal.NewMap()
		for key, value := range data {
			if _, ok := key.(string); !ok {
				// Cloudify DSL does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType("string")
			}
			name := yamlkeys.KeyString(key)
			map_.Put(name, NewValue(self.Context.MapChild(name, value)).Normalize())
		}
		constrainable = map_

	case *tosca.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		constrainable = normal.NewFunctionCall(data)

	default:
		value := normal.NewValue(data)
		value.Type = self.Type
		constrainable = value
	}

	if self.Description != nil {
		constrainable.SetDescription(*self.Description)
	}

	return constrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) SetIfNil(context *tosca.Context, key string, data interface{}) {
	if _, ok := self[key]; !ok {
		self[key] = NewValue(context.MapChild(key, data))
	}
}

func (self Values) RenderMissingValue(definition *ParameterDefinition, kind string, required bool, context *tosca.Context) {
	if definition.Default != nil {
		self[definition.Name] = definition.Default
	} else if required {
		context.MapChild(definition.Name, nil).ReportPropertyRequired(kind)
	} else if kind == "attribute" {
		// Parameters should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if value, ok := self[key]; !ok {
			// PropertyDefinition.Required defaults to true
			required := (definition.Required == nil) || *definition.Required
			self.RenderMissingValue(definition.ParameterDefinition, kind, required, context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}
}

func (self Values) RenderParameters(definitions ParameterDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if _, ok := self[key]; !ok {
			self.RenderMissingValue(definition, kind, false, context)
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared(kind)
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderParameter(definition.DataType, definition, true)
		}
	}
}

func (self Values) Normalize(c normal.Constrainables, prefix string) {
	for key, value := range self {
		c[prefix+key] = value.Normalize()
	}
}

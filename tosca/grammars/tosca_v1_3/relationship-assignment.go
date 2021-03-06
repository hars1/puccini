package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.2.2.3
//

type RelationshipAssignment struct {
	*Entity `name:"relationship"`

	RelationshipTemplateNameOrTypeName *string              `read:"type"`
	Properties                         Values               `read:"properties,Value"`
	Attributes                         Values               `read:"attributes,AttributeValue"` // missing in spec
	Interfaces                         InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	RelationshipTemplate *RelationshipTemplate `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
	RelationshipType     *RelationshipType     `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
}

func NewRelationshipAssignment(context *tosca.Context) *RelationshipAssignment {
	return &RelationshipAssignment{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Attributes: make(Values),
		Interfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadRelationshipAssignment(context *tosca.Context) interface{} {
	self := NewRelationshipAssignment(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.RelationshipTemplateNameOrTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

func (self *RelationshipAssignment) GetType() *RelationshipType {
	if self.RelationshipTemplate != nil {
		return self.RelationshipTemplate.RelationshipType
	} else if self.RelationshipType != nil {
		return self.RelationshipType
	} else {
		return nil
	}
}

func (self *RelationshipAssignment) Render(definition *RelationshipDefinition) {
	if definition != nil {
		// We will consider the "interfaces" at the definition to take priority over those at the type
		self.Interfaces.Render(definition.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))

		// Validate type compatibility
		relationshipType := self.GetType()
		if (relationshipType != nil) && (definition.RelationshipType != nil) && !self.Context.Hierarchy.IsCompatible(definition.RelationshipType, relationshipType) {
			self.Context.ReportIncompatibleType(relationshipType, definition.RelationshipType)
			return
		}
	}

	if self.RelationshipType != nil {
		self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
		self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
		self.Interfaces.Render(self.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
	} else if self.RelationshipTemplate != nil {
		self.Properties.CopyUnassigned(self.RelationshipTemplate.Properties)
		self.Attributes.CopyUnassigned(self.RelationshipTemplate.Attributes)
		self.Interfaces.CopyUnassigned(self.RelationshipTemplate.Interfaces)
	}
}

func (self *RelationshipAssignment) Normalize(r *normal.Relationship) {
	relationshipType := self.GetType()
	if (self.RelationshipTemplate != nil) && (self.RelationshipTemplate.Description != nil) {
		r.Description = *self.RelationshipTemplate.Description
	} else if (relationshipType != nil) && (relationshipType.Description != nil) {
		r.Description = *relationshipType.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, relationshipType); ok {
		r.Types = types
	}

	self.Properties.Normalize(r.Properties)
	self.Attributes.Normalize(r.Attributes)
	self.Interfaces.NormalizeForRelationship(self, r)
}

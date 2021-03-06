package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// SubstitutionMappings
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["AttributeMappings"] = ""
	context.ReadOverrides["SubstitutionFilter"] = ""

	return tosca_v1_3.ReadSubstitutionMappings(context)
}

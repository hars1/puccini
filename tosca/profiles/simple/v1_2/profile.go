// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/profile.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_2

metadata:
  puccini.scriptlet.import|tosca.lib.utils: internal:/tosca/common/1.0/js/lib/utils.js
  puccini.scriptlet.import|tosca.lib.coerce: internal:/tosca/common/1.0/js/lib/coerce.js
  puccini.scriptlet.import|tosca.coerce: internal:/tosca/common/1.0/js/coerce.js
  puccini.scriptlet.import|tosca.resolve: internal:/tosca/common/1.0/js/resolve.js

imports:
- artifacts.yaml
- groups.yaml
- nodes.yaml
- policies.yaml
`
}

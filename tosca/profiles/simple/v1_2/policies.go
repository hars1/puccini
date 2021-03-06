// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/policies.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_2

policy_types:

  tosca.policies.Root:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      citation_location: 5.11.1
    description: >-
      This is the default (root) TOSCA Policy Type definition that all other TOSCA base Policy Types
      derive from.

  tosca.policies.Placement:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      citation_location: 5.11.2
    description: >-
      This is the default (root) TOSCA Policy Type definition that is used to govern placement of
      TOSCA nodes or groups of nodes.
    derived_from: tosca.policies.Root

  tosca.policies.Scaling:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      citation_location: 5.11.3
    description: >-
      This is the default (root) TOSCA Policy Type definition that is used to govern scaling of
      TOSCA nodes or groups of nodes.
    derived_from: tosca.policies.Root

  tosca.policies.Update:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      citation_location: 5.11.4
    description: >-
      This is the default (root) TOSCA Policy Type definition that is used to govern update of TOSCA
      nodes or groups of nodes.
    derived_from: tosca.policies.Root

  tosca.policies.Performance:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      citation_location: 5.11.5
    description: >-
      This is the default (root) TOSCA Policy Type definition that is used to declare performance
      requirements for TOSCA nodes or groups of nodes.
    derived_from: tosca.policies.Root
`
}

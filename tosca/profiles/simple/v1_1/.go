// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/"] = `
# Modified from a file that was distributed with this NOTICE:
#
#   Apache AriaTosca
#   Copyright 2016-2017 The Apache Software Foundation
#
#   This product includes software developed at
#   The Apache Software Foundation (http://www.apache.org/).

tosca_definitions_version: tosca_simple_yaml_1_1

data_types:

  #
  # Primitive
  #

  string:
    metadata:
      puccini-tosca.type: string

  integer:
    metadata:
      puccini-tosca.type: integer

  float:
    metadata:
      puccini-tosca.type: float

  boolean:
    metadata:
      puccini-tosca.type: boolean

  timestamp:
    metadata:
      puccini-tosca.type: timestamp

  #
  # With entry schema
  #

  list:
    metadata:
      puccini-tosca.type: list
      specification: tosca-simple-1.1
      specification_section: 3.2.4
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_LIST'

  map:
    metadata:
      puccini-tosca.type: map
      specification: tosca-simple-1.1
      specification_section: 3.2.5
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_MAP'

  #
  # Special
  #

  version:
    metadata:
      puccini-tosca.type: version
      specification: tosca-simple-1.1
      specification_section: 3.2.2
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_VERSION'

  range:
    metadata:
      puccini-tosca.type: range
      specification: tosca-simple-1.1
      specification_section: 3.2.3
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_RANGE'

  #
  # Scalar
  #

  scalar-unit.size:
    metadata:
      puccini-tosca.type: scalar-unit.size
      specification: tosca-simple-1.1
      specification_section: 3.2.6.4
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_SCALAR_UNIT_SIZE'

  scalar-unit.time:
    metadata:
      puccini-tosca.type: scalar-unit.time
      specification: tosca-simple-1.1
      specification_section: 3.2.6.5
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_SCALAR_UNIT_TIME'

  scalar-unit.frequency:
    metadata:
      puccini-tosca.type: scalar-unit.frequency
      specification: tosca-simple-1.1
      specification_section: 3.2.6.6
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_SCALAR_UNIT_FREQUENCY'

  #
  # Complex
  #

  tosca.datatypes.Root:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.1
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_ROOT'
    description: >-
      This is the default (root) TOSCA Root Type definition that all complex TOSCA Data Types derive from.

  tosca.datatypes.Credential:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.2
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_CREDENTIAL'
    description: >-
      The Credential type is a complex TOSCA data Type used when describing authorization credentials used to access network
      accessible resources.
    derived_from: tosca.datatypes.Root
    properties:
      protocol:
        description: >-
          The optional protocol name.
        type: string
        required: false
      token_type:
        description: >-
          The required token type.
        type: string
        default: password
      token:
        description: >-
          The required token used as a credential for authorization or access to a networked resource.
        type: string
        required: false
      keys:
        description: >-
          The optional list of protocol-specific keys or assertions.
        type: map
        entry_schema:
          type: string
        required: false
      user:
        description: >-
          The optional user (name or ID) used for non-token based credentials.
        type: string
        required: false

  tosca.datatypes.network.NetworkInfo:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.3
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_NETWORKINFO'
    description: >-
      The Network type is a complex TOSCA data type used to describe logical network information.
    derived_from: tosca.datatypes.Root
    properties:
      network_name:
        description: >-
          The name of the logical network. e.g., "public", "private", "admin". etc.
        type: string
        required: false
      network_id:
        description: >-
          The unique ID of for the network generated by the network provider.
        type: string
        required: false
      addresses:
        description: >-
          The list of IP addresses assigned from the underlying network.
        type: list
        entry_schema:
          type: string
        required: false

  tosca.datatypes.network.PortInfo:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.4
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_PORTINFO'
    description: >-
      The PortInfo type is a complex TOSCA data type used to describe network port information.
    derived_from: tosca.datatypes.Root
    properties:
      port_name:
        description: >-
          The logical network port name.
        type: string
        required: false
      port_id:
        description: >-
          The unique ID for the network port generated by the network provider.
        type: string
        required: false
      network_id:
        description: >-
          The unique ID for the network.
        type: string
        required: false
      mac_address:
        description: >-
          The unique media access control address (MAC address) assigned to the port.
        type: string
        required: false
      addresses:
        description: >-
          The list of IP address(es) assigned to the port.
        type: list
        entry_schema:
          type: string
        required: false

  tosca.datatypes.network.PortDef:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.5
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_PORTDEF'
    description: >-
      The PortDef type is a TOSCA data Type used to define a network port.
    derived_from: integer # ARIA NOTE: we allow deriving from primitives
    constraints:
    - in_range: [ 1, 65535 ]

  tosca.datatypes.network.PortSpec:
    metadata:
      normative: 'true'
      specification: tosca-simple-1.1
      specification_section: 5.2.6
      specification_url: 'http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#TYPE_TOSCA_DATA_PORTSPEC'
    description: >-
      The PortSpec type is a complex TOSCA data Type used when describing port specifications for a network connection.
    derived_from: tosca.datatypes.Root
    properties:
      protocol:
        description: >-
          The required protocol used on the port.
        type: string
        constraints:
        - valid_values: [ udp, tcp, igmp ]
        default: tcp
      source:
        description: >-
          The optional source port.
        type: tosca.datatypes.network.PortDef
        required: false
      source_range:
        description: >-
          The optional range for source port.
        type: range
        constraints:
        - in_range: [ 1, 65535 ]
        required: false
      target:
        description: >-
          The optional target port.
        type: tosca.datatypes.network.PortDef
        required: false
      target_range:
        description: >-
          The optional range for target port.
        type: range
        constraints:
        - in_range: [ 1, 65535 ]
        required: false
`
}

package compiler

import (
	"github.com/tliron/puccini/ard"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/tosca/normal"
)

func Compile(s *normal.ServiceTemplate) (*cloutpkg.Clout, error) {
	clout := cloutpkg.NewClout()

	metadata := make(ard.StringMap)
	for name, scriptlet := range s.ScriptletNamespace {
		scriptlet, err := scriptlet.Read()
		if err != nil {
			return nil, err
		}
		if err = ard.StringMapPutNested(metadata, name, scriptlet); err != nil {
			return nil, err
		}
	}
	clout.Metadata["puccini-js"] = metadata

	history := ard.List{ard.StringMap{
		"timestamp":   common.Timestamp(),
		"description": "compile",
	}}

	metadata = make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["history"] = history
	clout.Metadata["puccini-tosca"] = metadata

	tosca := make(ard.StringMap)
	tosca["description"] = s.Description
	tosca["metadata"] = s.Metadata
	tosca["inputs"] = s.Inputs
	tosca["outputs"] = s.Outputs
	clout.Properties["tosca"] = tosca

	nodeTemplates := make(map[string]*cloutpkg.Vertex)

	// Node templates
	for _, nodeTemplate := range s.NodeTemplates {
		v := clout.NewVertex(cloutpkg.NewKey())

		nodeTemplates[nodeTemplate.Name] = v

		SetMetadata(v, "NodeTemplate")
		v.Properties["name"] = nodeTemplate.Name
		v.Properties["description"] = nodeTemplate.Description
		v.Properties["types"] = nodeTemplate.Types
		v.Properties["directives"] = nodeTemplate.Directives
		v.Properties["properties"] = nodeTemplate.Properties
		v.Properties["attributes"] = nodeTemplate.Attributes
		v.Properties["requirements"] = nodeTemplate.Requirements
		v.Properties["capabilities"] = nodeTemplate.Capabilities
		v.Properties["interfaces"] = nodeTemplate.Interfaces
		v.Properties["artifacts"] = nodeTemplate.Artifacts
	}

	groups := make(map[string]*cloutpkg.Vertex)

	// Groups
	for _, group := range s.Groups {
		v := clout.NewVertex(cloutpkg.NewKey())

		groups[group.Name] = v

		SetMetadata(v, "Group")
		v.Properties["name"] = group.Name
		v.Properties["description"] = group.Description
		v.Properties["types"] = group.Types
		v.Properties["properties"] = group.Properties
		v.Properties["interfaces"] = group.Interfaces

		for _, nodeTemplate := range group.Members {
			nv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(nv)

			SetMetadata(e, "Member")
		}
	}

	workflows := make(map[string]*cloutpkg.Vertex)

	// Workflows
	for _, workflow := range s.Workflows {
		v := clout.NewVertex(cloutpkg.NewKey())

		workflows[workflow.Name] = v

		SetMetadata(v, "Workflow")
		v.Properties["name"] = workflow.Name
		v.Properties["description"] = workflow.Description
	}

	// Workflow steps
	for name, workflow := range s.Workflows {
		v := workflows[name]

		steps := make(map[string]*cloutpkg.Vertex)

		for _, step := range workflow.Steps {
			sv := clout.NewVertex(cloutpkg.NewKey())

			steps[step.Name] = sv

			SetMetadata(sv, "WorkflowStep")
			sv.Properties["name"] = step.Name

			e := v.NewEdgeTo(sv)
			SetMetadata(e, "WorkflowStep")

			if step.TargetNodeTemplate != nil {
				nv := nodeTemplates[step.TargetNodeTemplate.Name]
				e = sv.NewEdgeTo(nv)
				SetMetadata(e, "NodeTemplateTarget")
			} else if step.TargetGroup != nil {
				gv := groups[step.TargetGroup.Name]
				e = sv.NewEdgeTo(gv)
				SetMetadata(e, "GroupTarget")
			} else {
				// This would happen only if there was a parsing error
				continue
			}

			// Workflow activities
			for sequence, activity := range step.Activities {
				av := clout.NewVertex(cloutpkg.NewKey())

				e = sv.NewEdgeTo(av)
				SetMetadata(e, "WorkflowActivity")
				e.Properties["sequence"] = sequence

				SetMetadata(av, "WorkflowActivity")
				if activity.DelegateWorkflow != nil {
					wv := workflows[activity.DelegateWorkflow.Name]
					e = av.NewEdgeTo(wv)
					SetMetadata(e, "DelegateWorkflow")
				} else if activity.InlineWorkflow != nil {
					wv := workflows[activity.InlineWorkflow.Name]
					e = av.NewEdgeTo(wv)
					SetMetadata(e, "InlineWorkflow")
				} else if activity.SetNodeState != "" {
					av.Properties["setNodeState"] = activity.SetNodeState
				} else if activity.CallOperation != nil {
					m := make(ard.StringMap)
					m["interface"] = activity.CallOperation.Interface.Name
					m["operation"] = activity.CallOperation.Name
					av.Properties["callOperation"] = m
				}
			}
		}

		for _, step := range workflow.Steps {
			sv := steps[step.Name]

			for _, next := range step.OnSuccessSteps {
				nsv := steps[next.Name]
				e := sv.NewEdgeTo(nsv)
				SetMetadata(e, "OnSuccess")
			}

			for _, next := range step.OnFailureSteps {
				nsv := steps[next.Name]
				e := sv.NewEdgeTo(nsv)
				SetMetadata(e, "OnFailure")
			}
		}
	}

	// Policies
	for _, policy := range s.Policies {
		v := clout.NewVertex(cloutpkg.NewKey())

		SetMetadata(v, "Policy")
		v.Properties["name"] = policy.Name
		v.Properties["description"] = policy.Description
		v.Properties["types"] = policy.Types
		v.Properties["properties"] = policy.Properties

		for _, nodeTemplate := range policy.NodeTemplateTargets {
			nv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(nv)

			SetMetadata(e, "NodeTemplateTarget")
		}

		for _, group := range policy.GroupTargets {
			gv := groups[group.Name]
			e := v.NewEdgeTo(gv)

			SetMetadata(e, "GroupTarget")
		}

		for _, trigger := range policy.Triggers {
			if trigger.Operation != nil {
				to := clout.NewVertex(cloutpkg.NewKey())

				SetMetadata(to, "Operation")
				to.Properties["description"] = trigger.Operation.Description
				to.Properties["implementation"] = trigger.Operation.Implementation
				to.Properties["dependencies"] = trigger.Operation.Dependencies
				to.Properties["inputs"] = trigger.Operation.Inputs

				e := v.NewEdgeTo(to)
				SetMetadata(e, "PolicyTriggerOperation")
			} else if trigger.Workflow != nil {
				wv := workflows[trigger.Workflow.Name]

				e := v.NewEdgeTo(wv)
				SetMetadata(e, "PolicyTriggerWorkflow")
			}
		}
	}

	// Substitution
	if s.Substitution != nil {
		v := clout.NewVertex(cloutpkg.NewKey())

		SetMetadata(v, "Substitution")
		v.Properties["type"] = s.Substitution.Type
		v.Properties["typeMetadata"] = s.Substitution.TypeMetadata

		for nodeTemplate, capability := range s.Substitution.CapabilityMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "CapabilityMapping")
			e.Properties["capability"] = capability.Name
		}

		for nodeTemplate, requirement := range s.Substitution.RequirementMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "RequirementMapping")
			e.Properties["requirement"] = requirement
		}

		for nodeTemplate, property := range s.Substitution.PropertyMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "PropertyMapping")
			e.Properties["property"] = property
		}

		for nodeTemplate, attribute := range s.Substitution.AttributeMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "AttributeMapping")
			e.Properties["attribute"] = attribute
		}

		for nodeTemplate, interface_ := range s.Substitution.InterfaceMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "InterfaceMapping")
			e.Properties["interface"] = interface_
		}
	}

	// Normalize
	var err error
	clout, err = clout.Normalize()
	if err != nil {
		return clout, err
	}

	// TODO: call JavaScript plugins

	return clout, nil
}

func SetMetadata(entity cloutpkg.Entity, kind string) {
	metadata := make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["kind"] = kind
	entity.GetMetadata()["puccini-tosca"] = metadata
}

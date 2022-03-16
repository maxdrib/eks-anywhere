package cluster

import anywherev1 "github.com/aws/eks-anywhere/pkg/api/v1alpha1"

func awsIamEntry() *ConfigManagerEntry {
	return &ConfigManagerEntry{
		APIObjectMapping: map[string]APIObjectGenerator{
			anywherev1.AWSIamConfigKind: func() APIObject {
				return &anywherev1.AWSIamConfig{}
			},
		},
		Processors: []ParsedProcessor{processAWSIam},
	}
}

func processAWSIam(c *Config, objects ObjectLookup) {
	if c.AWSIAMConfigs == nil {
		c.AWSIAMConfigs = map[string]*anywherev1.AWSIamConfig{}
	}

	for _, idr := range c.Cluster.Spec.IdentityProviderRefs {
		idp := objects.GetFromRef(c.Cluster.APIVersion, idr)
		if idp == nil {
			return
		}
		if idr.Kind == anywherev1.AWSIamConfigKind {
			c.AWSIAMConfigs[idp.GetName()] = idp.(*anywherev1.AWSIamConfig)
		}
	}
}

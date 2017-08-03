package cluster

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/patrickmn/go-cache"
	"github.com/stitchfix/flotilla-os/config"
	"github.com/stitchfix/flotilla-os/state"
	"strings"
	"time"
)

//
// ECSClusterClient is the default cluster client and maintains a
// cached map[string]instanceResources which is used to check that
// the resources requested by a definition  -at some point- could
// become available on the specified cluster.
//
// [NOTE] This client assumes homogenous clusters
//
type ECSClusterClient struct {
	ecsClient resourceClient
	clusters  resourceCache
}

type resourceClient interface {
	DescribeClusters(input *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error)
	ListContainerInstances(input *ecs.ListContainerInstancesInput) (*ecs.ListContainerInstancesOutput, error)
	DescribeContainerInstances(input *ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error)
}

type instanceResources struct {
	memory int64
	cpu    int64
}

//
// Name is the name of the client
//
func (ecc *ECSClusterClient) Name() string {
	return "ecs"
}

//
// Initialize the ecs cluster client with config
//
func (ecc *ECSClusterClient) Initialize(conf config.Config) error {
	if !conf.IsSet("aws_default_region") {
		return fmt.Errorf("ecsClusterClient needs [aws_default_region] set in config")
	}

	flotillaMode := conf.GetString("flotilla_mode")
	if flotillaMode != "test" {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(conf.GetString("aws_default_region"))}))

		ecc.ecsClient = ecs.New(sess)
	}
	ecc.clusters = resourceCache{
		duration:      15 * time.Minute,
		internalCache: cache.New(15*time.Minute, 5*time.Minute),
	}
	return nil
}

//
// CanBeRun determines whether a task formed from the specified definition
// can be run on clusterName
//
func (ecc *ECSClusterClient) CanBeRun(clusterName string, definition state.Definition) (bool, error) {
	var (
		resources *instanceResources
		err       error
	)
	resources, found := ecc.clusters.getInstanceResources(clusterName)
	if !found {
		resources, err = ecc.fetchResources(clusterName)
		if err != nil {
			return false, err
		}
		if resources == nil {
			return false, nil
		}
		ecc.clusters.setInstanceResources(clusterName, resources)
	}

	return ecc.validate(resources, definition), nil
}

func (ecc *ECSClusterClient) validate(resources *instanceResources, definition state.Definition) bool {
	if resources != nil && definition.Memory != nil && int64(*definition.Memory) < resources.memory {
		// TODO - check cpu when available on the definition
		return true
	}
	return false
}

func (ecc *ECSClusterClient) fetchResources(clusterName string) (*instanceResources, error) {
	exists, err := ecc.clusterExists(clusterName)
	if err != nil {
		return nil, err
	}
	if exists {
		return ecc.clusterInstanceResources(clusterName)
	}
	return nil, nil
}

func (ecc *ECSClusterClient) clusterExists(clusterName string) (bool, error) {
	result, err := ecc.ecsClient.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{
			&clusterName,
		},
	})

	if err != nil {
		return false, err
	}
	if len(result.Failures) != 0 {
		msg := make([]string, len(result.Failures))
		for i, failure := range result.Failures {
			msg[i] = *failure.Reason
		}
		return false, fmt.Errorf("ERRORS: %s", strings.Join(msg, "\n"))
	}
	if len(result.Clusters) == 0 {
		return false, nil
	}
	return true, nil
}

func (ecc *ECSClusterClient) clusterInstanceResources(clusterName string) (*instanceResources, error) {
	var result instanceResources

	instances, err := ecc.listInstances(&ecs.ListContainerInstancesInput{
		Cluster: &clusterName,
	})
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, nil // short-circuit to avoid additional spurious api call with zero instances
	}

	resources, err := ecc.describeInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            &clusterName,
		ContainerInstances: instances,
	})

	if err != nil {
		return nil, err
	}

	//
	// Assumes -all- instances in the cluster have identical registered resources
	// While this is certainly the most manageable way to configure ecs clusters
	// it is -not- the only way in practice
	//
	if resources != nil && len(resources) > 0 {
		// An alternative, and potentially better handling of heterogeneous clusters,
		// is to return the resources with the -highest- memory and cpu
		result = resources[0]
	}
	return &result, nil
}

func (ecc *ECSClusterClient) listInstances(input *ecs.ListContainerInstancesInput) ([]*string, error) {
	result, err := ecc.ecsClient.ListContainerInstances(input)
	if err != nil {
		return nil, err
	}
	var subset []*string
	for _, arn := range result.ContainerInstanceArns {
		subset = append(subset, arn)
	}

	if result.NextToken != nil {
		input.NextToken = result.NextToken
		more, err := ecc.listInstances(input)
		if err != nil {
			return nil, err
		}
		subset = append(subset, more...)
	}
	return subset, nil
}

func (ecc *ECSClusterClient) describeInstances(input *ecs.DescribeContainerInstancesInput) ([]instanceResources, error) {
	result, err := ecc.ecsClient.DescribeContainerInstances(input)
	if err != nil {
		return nil, err
	}

	if len(result.Failures) != 0 {
		msg := make([]string, len(result.Failures))
		for i, failure := range result.Failures {
			msg[i] = *failure.Reason
		}
		return nil, fmt.Errorf("ERRORS: %s", strings.Join(msg, "\n"))
	}

	res := make([]instanceResources, len(result.ContainerInstances))
	for i, ci := range result.ContainerInstances {
		irs := instanceResources{}
		for _, rsrc := range ci.RegisteredResources {
			if *rsrc.Name == "CPU" {
				irs.cpu = *rsrc.IntegerValue
			} else if *rsrc.Name == "MEMORY" {
				irs.memory = *rsrc.IntegerValue
			}
		}
		res[i] = irs
	}
	return res, nil
}

type resourceCache struct {
	duration      time.Duration
	internalCache *cache.Cache
}

func (rc *resourceCache) getInstanceResources(clusterName string) (*instanceResources, bool) {
	resources, found := rc.internalCache.Get(clusterName)
	if found {
		ir := resources.(instanceResources)
		return &ir, true
	}
	return &instanceResources{}, false
}

func (rc *resourceCache) setInstanceResources(clusterName string, resources *instanceResources) {
	rc.internalCache.Set(clusterName, *resources, rc.duration)
}
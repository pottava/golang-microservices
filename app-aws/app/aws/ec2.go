package aws

/**
 * @see https://github.com/aws/aws-sdk-go/blob/master/service/ec2/api.go
 */
import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pottava/golang-microservices/app-aws/app/logs"
)

// Ec2Instance returns a specified ec2 instance
func Ec2Instance(id string) (instance *ec2.Instance, e error) {
	req := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{awssdk.String(id)},
	}
	res, err := ec2.New(session.New(), awssdk.NewConfig()).DescribeInstances(req)
	if err != nil {
		return nil, err
	}
	for idx := range res.Reservations {
		for _, inst := range res.Reservations[idx].Instances {
			if *inst.InstanceId == id {
				instance = inst
			}
		}
	}
	return instance, nil
}

// Ec2Instances responses ec2 instances
func Ec2Instances() (instances []ec2.Instance, e error) {
	res, err := ec2.New(session.New(), awssdk.NewConfig()).DescribeInstances(nil)
	if err != nil {
		logs.Error.Print("Could not describe EC2 Instances.")
		return nil, err
	}
	instances = make([]ec2.Instance, len(res.Reservations))
	for idx := range res.Reservations {
		for _, inst := range res.Reservations[idx].Instances {
			instances[idx] = *inst
		}
	}
	return instances, nil
}

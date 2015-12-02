package models

// EC2Instance represents ec2 instance
type EC2Instance struct {
	InstanceID string `json:"InstanceId"`
}

type daoEC2Instance struct {
	Header   APIHeader      `json:"header"`
	Response []*EC2Instance `json:"response"`
}

// GetEC2Instances retrives ec2 instances
func GetEC2Instances() (instances []*EC2Instance, found bool) {
	res := &daoEC2Instance{}
	if err := aws("GET", "/ec2/instances/", "", res); err == nil {
		if res.Header.Status == "success" {
			return res.Response, true
		}
	}
	return []*EC2Instance{}, false
}

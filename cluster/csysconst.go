package cluster

import (
	"errors"
)


var errJobConsumerNotRegister = errors.New("job consumer not register")
var errjobNotConfig =errors.New("job not config")

const (
	dsbServerRoot        = "@domain/dsbcenter/servers"
	dsbServerPath        = "@domain/dsbcenter/servers/dsb_"
	dsbServerValue       = `{"domain":"@domain","path":"@path","ip":"@ip","port":"@port",server":"@type","online":"@online","lastPublish":"@pst","last":"@last"}`
	dsbServerConfig      = "@domain/configs/dsbcenter/config"
	servicePublishPath   = "@domain/configs/services/publish"
	serviceRoot          = "@domain/services"
	serviceProviderRoot  = "@domain/services/@serviceName/providers"
	serviceProviderPath  = "@domain/services/@serviceName/providers/@ip"
	serviceProviderValue = `{"last":@now}`

	serviceConsumerRoot  = "@domain/services/@serviceName/consumers"
	serviceConsumerPath  = "@domain/services/@serviceName/consumers/@ip"
	serviceConsumerValue = `{"last":@now}`
	jobRoot              = "@domain/jobs"
	jobConfigPath        = "@domain/configs/jobs/config"
	jobConsumerRoot      = "@domain/jobs/@jobName"
	jobConsumerPath      = "@domain/jobs/@jobName/job"
	jobConsumerValue     = `{"ip":"@ip",last":@now}`
)

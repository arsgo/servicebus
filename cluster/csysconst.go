package cluster

import (
	"errors"
)

var errJobConsumerNotRegister = errors.New("job consumer not register")
var errjobNotConfig = errors.New("job not config")

const (
	dsbServerRoot        = "@domain/rc/servers"
	dsbServerPath        = "@domain/rc/servers/rc_"
	dsbServerValue       = `{"domain":"@domain","path":"@path","ip":"@ip","port":"@port",server":"@type","online":"@online","lastPublish":"@pst","last":"@last"}`
	dsbServerConfig      = "@domain/configs/rc/config"
	servicePublishPath   = "@domain/configs/sp/publish"
	serviceRoot          = "@domain/sp"
	serviceConfig        = "@domain/configs/sp/config"
	serviceProviderRoot  = "@domain/sp/@serviceName/providers"
	serviceProviderPath  = "@domain/sp/@serviceName/providers/@ip@port"
	serviceProviderValue = `{"last":@last}`

	serviceConsumerRoot  = "@domain/sp/@serviceName/consumers"
	serviceConsumerPath  = "@domain/sp/@serviceName/consumers/@ip"
	serviceConsumerValue = `{"last":@now}`
	jobRoot              = "@domain/job"
	jobConfigPath        = "@domain/configs/job/config"
	jobConsumerRoot      = "@domain/job/@jobName"
	jobConsumerPath      = "@domain/job/@jobName/job"
	jobConsumerValue     = `{"ip":"@ip",last":@now}`

	appServerConfig   = "@domain/configs/app/@ip"
	appServerRoot     = "@domain/app/servers"
	appServerRootPath = "@domain/app/servers/@ip"
)

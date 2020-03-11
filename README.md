# metrics-beaver
Fetches metrics from cloudwatch. It can be used with a monitoring system like nagios or zabbix to fetch metrics from cloudwatch. it's installed in the zabbix server in this [script](https://github.com/spring-media/la-hollywood/blob/master/zabbix_bootstrap/beaver_bootstrap.sh#L89).

#### Installing
	go get github.com/spring-media/la-metrics-beaver 
	cd $GOPATH/github/spring-media/la-metrics-beaver
	go build main.go

#### Examples
	
	./main  \
	  -namespace AWS/ApplicationELB \
	  -metric.name RequestCount \
	  -dimension.name LoadBalancer \
	  -dimension.value app/ECS-la-production-ecs/1faecb4f1f100e1b \
	  -stats Sum \
	  -period 5m

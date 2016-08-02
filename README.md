# metrics-beaver
Fetches metrics from cloudwatch. It can be used with a monitoring system like nagios or zabbix to fetch metrics from cloudwath.

#### Installing
        go get github.com/mrsn/metrics-beaver 

#### Examples
	
	metrics-beaver -namespace AWS/ELB \
	-metric-name RequestCount \
	-monitoring-type detailed \
	-statistics Sum \
	-dimension-name LoadBalancerName \
	-dimension-value elb-foobar

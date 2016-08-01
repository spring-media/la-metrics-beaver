# metric-beaver
fetches metrics from cloudwatch

Examples:
	./main -namespace AWS/ELB -metric-name RequestCount -monitoring-type detailed --statistics Sum -dimension-name LoadBalancerName -dimension-value elb-foobar


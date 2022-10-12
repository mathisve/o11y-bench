#!/bin/bash

hey -q 1 -c 100 -n 10000 \
    http://ac5a5f124960b41a28a5f487ea52bb27-223003468.us-east-1.elb.amazonaws.com:16686/api/traces\?end\=1665605526162000\&limit\=20\&lookback\=5m\&maxDuration\=\&minDuration\=\&operation\=lets-go\&service\=tracegen\&start\=1665605226162000
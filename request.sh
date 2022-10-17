#!/bin/bash
echo "testing"
hey -q 1 -c 100 -n 100 	http://ac5a5f124960b41a28a5f487ea52bb27-529402917.us-east-1.elb.amazonaws.com:16686/api/traces?end=1665675851998000&limit=20&lookback=3h&maxDuration=&minDuration=&service=tracegen&start=1665665051998000
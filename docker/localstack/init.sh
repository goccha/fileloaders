#!/bin/sh
awslocal ssm put-parameter --name '/parameter/test' --type 'String' --value 'value/fileloaders/test'
awslocal ssm put-parameter --name '/parameter/test1' --type 'String' --value 'value/fileloaders/test1'
awslocal ssm put-parameter --name '/parameter1/test' --type 'String' --value 'value/fileloaders/test'
awslocal ssm put-parameter --name '/parameter1/secure' --type 'SecureString' --value 'secure-value/fileloaders/test'

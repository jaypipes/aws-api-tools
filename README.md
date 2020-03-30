# `aws-api-tools`

Tools to explore and transform AWS service APIs.

## `aws-api-tool` CLI binary

The `aws-api-tool` CLI binary can be used to list AWS service APIs and explore
aspects of those APIs.

### Installation

The easiest way to install `aws-api-tool` is to use `go get`:

```
$ go get github.com/jaypipes/aws-api-tools/cmd/aws-api-tool
```

which will install the binary into your `$GOPATH/bin` directory.

### List AWS service APIs

Use the `aws-api-tool service list` command to list AWS services:

```
$ aws-api-tool service list
+------------------------------+-------------+------------------------------------+
|            ALIAS             | API VERSION |             FULL NAME              |
+------------------------------+-------------+------------------------------------+
| AWSMigrationHub              | 2017-05-31  | AWS Migration Hub                  |
| accessanalyzer               | 2019-11-01  | Access Analyzer                    |
<snip>
| workspaces                   | 2015-04-08  | Amazon WorkSpaces                  |
| xray                         | 2016-04-12  | AWS X-Ray                          |
+------------------------------+-------------+------------------------------------+
```

You can filter the results using the `--filter` flag, which accepts a
comma-delimited string of strings to match for the service's alias:

```
$ aws-api-tool service list --filter ebs,sdb
+-------+-------------+----------------------------+
| ALIAS | API VERSION |         FULL NAME          |
+-------+-------------+----------------------------+
| ebs   | 2019-11-02  | Amazon Elastic Block Store |
| sdb   | 2009-04-15  | Amazon SimpleDB            |
+-------+-------------+----------------------------+
```

### Explore an AWS service API

Explore information about a particular AWS service API with the `aws-api-tool
api` command. This command and all of its subcommands accepts a single
`--service` argument which should be the alias of a service.

#### Get summary information about an API

To get summary information about a particular AWS service API, use the
`aws-api-tool api info` command.

```
$ aws-api-tool api --service sns info
Full name:        Amazon Simple Notification Service
API version:      2010-03-31
Total operations: 33
Total primaries:  3
Total objects:    6
Total scalars:    30
Total payloads:   56
Total exceptions: 23
Total lists:      10
```

#### List API operations

To list operations for an AWS service API, use the `aws-api-tool api
operations` command:

```
$ aws-api-tool api --service sns operations
+------------------------------------+-------------+
|                NAME                | HTTP METHOD |
+------------------------------------+-------------+
| AddPermission                      | POST        |
| CheckIfPhoneNumberIsOptedOut       | POST        |
| ConfirmSubscription                | POST        |
| CreatePlatformApplication          | POST        |
| CreatePlatformEndpoint             | POST        |
| CreateTopic                        | POST        |
| DeleteEndpoint                     | POST        |
| DeletePlatformApplication          | POST        |
| DeleteTopic                        | POST        |
| GetEndpointAttributes              | POST        |
| GetPlatformApplicationAttributes   | POST        |
| GetSMSAttributes                   | POST        |
| GetSubscriptionAttributes          | POST        |
| GetTopicAttributes                 | POST        |
| ListEndpointsByPlatformApplication | POST        |
| ListPhoneNumbersOptedOut           | POST        |
| ListPlatformApplications           | POST        |
| ListSubscriptions                  | POST        |
| ListSubscriptionsByTopic           | POST        |
| ListTagsForResource                | POST        |
| ListTopics                         | POST        |
| OptInPhoneNumber                   | POST        |
| Publish                            | POST        |
| RemovePermission                   | POST        |
| SetEndpointAttributes              | POST        |
| SetPlatformApplicationAttributes   | POST        |
| SetSMSAttributes                   | POST        |
| SetSubscriptionAttributes          | POST        |
| SetTopicAttributes                 | POST        |
| Subscribe                          | POST        |
| TagResource                        | POST        |
| Unsubscribe                        | POST        |
| UntagResource                      | POST        |
+------------------------------------+-------------+
```

#### List API primary objects

Primary objects are those objects that are "top-level" constructs in an API.
These primary objects correspond to the core structures exposed in the API with
Create, Read, Update and Delete operations.

Use the `aws-api-tool api primaries` command to list these primary objects:

```
$ aws-api-tool api --service sqs primaries
+-------+
| NAME  |
+-------+
| Queue |
+-------+
```

#### List API objects

Objects are those elements of the API that are **NOT** Lists, Payloads or
Scalar types. You can think of objects as complex types that contain a set of
fields that may be other objects or scalar types.

Use the `aws-api-tool api objects` command to list an API's objects:

```
$ aws-api-tool api --service ecr objects
+-------------------------------+
|             NAME              |
+-------------------------------+
| Attribute                     |
| AuthorizationData             |
| DescribeImagesFilter          |
| Image                         |
| ImageDetail                   |
| ImageFailure                  |
| ImageIdentifier               |
| ImageScanFinding              |
| ImageScanFindings             |
| ImageScanFindingsSummary      |
| ImageScanStatus               |
| ImageScanningConfiguration    |
| Layer                         |
| LayerFailure                  |
| LifecyclePolicyPreviewFilter  |
| LifecyclePolicyPreviewResult  |
| LifecyclePolicyPreviewSummary |
| LifecyclePolicyRuleAction     |
| ListImagesFilter              |
| Repository                    |
| Tag                           |
+-------------------------------+
```

#### List API scalars

Scalars are concepts of the API that are neither structure objects or list
objects. They are the simplest data types in the API.

Use the `aws-api-tool api scalars` command to list scalars:

```
$ go run cmd/aws-api-tool/main.go api --service sqs scalars
+------------------------------------+---------+
|                NAME                |  TYPE   |
+------------------------------------+---------+
| Binary                             | blob    |
| Boolean                            | boolean |
| Integer                            | integer |
| MessageAttributeName               | string  |
| MessageBodyAttributeMap            | map     |
| MessageBodySystemAttributeMap      | map     |
| MessageSystemAttributeMap          | map     |
| MessageSystemAttributeName         | string  |
| MessageSystemAttributeNameForSends | string  |
| QueueAttributeMap                  | map     |
| QueueAttributeName                 | string  |
| String                             | string  |
| TagKey                             | string  |
| TagMap                             | map     |
| TagValue                           | string  |
+------------------------------------+---------+
```

#### List API list objects

List objects are constructs in the API that serve only as collections of a
single scalar or object type.

List an API's list objects with the `aws-api-tool api lists` command:

```
$ aws-api-tool api --service s3 lists
+---------------------------------+
|              NAME               |
+---------------------------------+
| AllowedHeaders                  |
| AllowedMethods                  |
| AllowedOrigins                  |
| AnalyticsConfigurationList      |
| Buckets                         |
| CORSRules                       |
| CommonPrefixList                |
| CompletedPartList               |
| DeleteMarkers                   |
| DeletedObjects                  |
| Errors                          |
| EventList                       |
| ExposeHeaders                   |
| FilterRuleList                  |
| Grants                          |
| InventoryConfigurationList      |
| InventoryOptionalFields         |
| LambdaFunctionConfigurationList |
| LifecycleRules                  |
| MetricsConfigurationList        |
| MultipartUploadList             |
| NoncurrentVersionTransitionList |
| ObjectIdentifierList            |
| ObjectList                      |
| ObjectVersionList               |
| Parts                           |
| QueueConfigurationList          |
| ReplicationRules                |
| RoutingRules                    |
| Rules                           |
| ServerSideEncryptionRules       |
| TagSet                          |
| TargetGrants                    |
| TopicConfigurationList          |
| TransitionList                  |
| UserMetadata                    |
+---------------------------------+
```

#### List API payloads

Payload objects are constructs in the API that represent the input or output
payload for an individual HTTP request.

Use the `aws-api-tool api payloads` command to list these payload objects:


```
$ go run cmd/aws-api-tool/main.go api --service ecr payloads
+---------------------------------------+
|                 NAME                  |
+---------------------------------------+
| BatchCheckLayerAvailabilityRequest    |
| BatchCheckLayerAvailabilityResponse   |
| BatchDeleteImageRequest               |
| BatchDeleteImageResponse              |
| BatchGetImageRequest                  |
| BatchGetImageResponse                 |
| CompleteLayerUploadRequest            |
| CompleteLayerUploadResponse           |
| CreateRepositoryRequest               |
| CreateRepositoryResponse              |
| DeleteLifecyclePolicyRequest          |
| DeleteLifecyclePolicyResponse         |
| DeleteRepositoryPolicyRequest         |
| DeleteRepositoryPolicyResponse        |
| DeleteRepositoryRequest               |
| DeleteRepositoryResponse              |
| DescribeImageScanFindingsRequest      |
| DescribeImageScanFindingsResponse     |
| DescribeImagesRequest                 |
| DescribeImagesResponse                |
| DescribeRepositoriesRequest           |
| DescribeRepositoriesResponse          |
| GetAuthorizationTokenRequest          |
| GetAuthorizationTokenResponse         |
| GetDownloadUrlForLayerRequest         |
| GetDownloadUrlForLayerResponse        |
| GetLifecyclePolicyPreviewRequest      |
| GetLifecyclePolicyPreviewResponse     |
| GetLifecyclePolicyRequest             |
| GetLifecyclePolicyResponse            |
| GetRepositoryPolicyRequest            |
| GetRepositoryPolicyResponse           |
| InitiateLayerUploadRequest            |
| InitiateLayerUploadResponse           |
| ListImagesRequest                     |
| ListImagesResponse                    |
| ListTagsForResourceRequest            |
| ListTagsForResourceResponse           |
| PutImageRequest                       |
| PutImageResponse                      |
| PutImageScanningConfigurationRequest  |
| PutImageScanningConfigurationResponse |
| PutImageTagMutabilityRequest          |
| PutImageTagMutabilityResponse         |
| PutLifecyclePolicyRequest             |
| PutLifecyclePolicyResponse            |
| SetRepositoryPolicyRequest            |
| SetRepositoryPolicyResponse           |
| StartImageScanRequest                 |
| StartImageScanResponse                |
| StartLifecyclePolicyPreviewRequest    |
| StartLifecyclePolicyPreviewResponse   |
| TagResourceRequest                    |
| TagResourceResponse                   |
| UntagResourceRequest                  |
| UntagResourceResponse                 |
| UploadLayerPartRequest                |
| UploadLayerPartResponse               |
+---------------------------------------+
```

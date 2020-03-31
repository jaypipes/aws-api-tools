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
Total resources:  3
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

You can filter the results using the `--method` and `--prefix` flags which both
take a comma-delimited string of HTTP methods or string name prefixes to filter
results by:

```
$ aws-api-tool api --service ec2 operations --prefix Update,Delete
+--------------------------------------------+-------------+
|                    NAME                    | HTTP METHOD |
+--------------------------------------------+-------------+
| DeleteClientVpnEndpoint                    | POST        |
| DeleteClientVpnRoute                       | POST        |
| DeleteCustomerGateway                      | POST        |
| DeleteDhcpOptions                          | POST        |
| DeleteEgressOnlyInternetGateway            | POST        |
| DeleteFleets                               | POST        |
| DeleteFlowLogs                             | POST        |
| DeleteFpgaImage                            | POST        |
| DeleteInternetGateway                      | POST        |
| DeleteKeyPair                              | POST        |
| DeleteLaunchTemplate                       | POST        |
| DeleteLaunchTemplateVersions               | POST        |
| DeleteLocalGatewayRoute                    | POST        |
| DeleteLocalGatewayRouteTableVpcAssociation | POST        |
| DeleteNatGateway                           | POST        |
| DeleteNetworkAcl                           | POST        |
| DeleteNetworkAclEntry                      | POST        |
| DeleteNetworkInterface                     | POST        |
| DeleteNetworkInterfacePermission           | POST        |
| DeletePlacementGroup                       | POST        |
| DeleteQueuedReservedInstances              | POST        |
| DeleteRoute                                | POST        |
| DeleteRouteTable                           | POST        |
| DeleteSecurityGroup                        | POST        |
| DeleteSnapshot                             | POST        |
| DeleteSpotDatafeedSubscription             | POST        |
| DeleteSubnet                               | POST        |
| DeleteTags                                 | POST        |
| DeleteTrafficMirrorFilter                  | POST        |
| DeleteTrafficMirrorFilterRule              | POST        |
| DeleteTrafficMirrorSession                 | POST        |
| DeleteTrafficMirrorTarget                  | POST        |
| DeleteTransitGateway                       | POST        |
| DeleteTransitGatewayMulticastDomain        | POST        |
| DeleteTransitGatewayPeeringAttachment      | POST        |
| DeleteTransitGatewayRoute                  | POST        |
| DeleteTransitGatewayRouteTable             | POST        |
| DeleteTransitGatewayVpcAttachment          | POST        |
| DeleteVolume                               | POST        |
| DeleteVpc                                  | POST        |
| DeleteVpcEndpointConnectionNotifications   | POST        |
| DeleteVpcEndpointServiceConfigurations     | POST        |
| DeleteVpcEndpoints                         | POST        |
| DeleteVpcPeeringConnection                 | POST        |
| DeleteVpnConnection                        | POST        |
| DeleteVpnConnectionRoute                   | POST        |
| DeleteVpnGateway                           | POST        |
| UpdateSecurityGroupRuleDescriptionsEgress  | POST        |
| UpdateSecurityGroupRuleDescriptionsIngress | POST        |
+--------------------------------------------+-------------+
$ aws-api-tool api --service s3 operations --method GET --prefix ListO
+--------------------+-------------+
|        NAME        | HTTP METHOD |
+--------------------+-------------+
| ListObjectVersions | GET         |
| ListObjects        | GET         |
| ListObjectsV2      | GET         |
+--------------------+-------------+
```
#### List API resource objects

Resource objects are those objects that are "top-level" constructs in an API.
These resource objects correspond to the core structures exposed in the API with
Create, Read, Update and Delete operations.

Use the `aws-api-tool api resources` command to list these resource objects:

```
$ aws-api-tool api --service sqs resources
+-------+
| NAME  |
+-------+
| Queue |
+-------+
```

Resource objects are only top-level objects in the API. If an object is solely
contained within another object, it is not a resource object. For example, the
AWS APIGateway API has the following Create operations:

```
$ aws-api-tool api --service apigateway operations --prefix Create
+----------------------------+-------------+
|            NAME            | HTTP METHOD |
+----------------------------+-------------+
| CreateApiKey               | POST        |
| CreateAuthorizer           | POST        |
| CreateBasePathMapping      | POST        |
| CreateDeployment           | POST        |
| CreateDocumentationPart    | POST        |
| CreateDocumentationVersion | POST        |
| CreateDomainName           | POST        |
| CreateModel                | POST        |
| CreateRequestValidator     | POST        |
| CreateResource             | POST        |
| CreateRestApi              | POST        |
| CreateStage                | POST        |
| CreateUsagePlan            | POST        |
| CreateUsagePlanKey         | POST        |
| CreateVpcLink              | POST        |
+----------------------------+-------------+
```

However, of the above, only the `ApiKey`, `DomainName`, `RestApi`, `UsagePlan` and
`VpcLink` are resources:

```
$ aws-api-tool api --service apigateway resources
+------------+
|    NAME    |
+------------+
| ApiKey     |
| DomainName |
| RestApi    |
| UsagePlan  |
| VpcLink    |
+------------+
```

This is because the other objects are solely contained within
another object. For example, a `Deployment` is solely a part of a `RestApi`
object; it cannot be created as a separate thing.

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

#### Show OpenAPI3 Schema for API resource

Use the `aws-api-tool schema openapi` command to display the OpenAPI3 Schema
for a specific resource in an AWS Service API. Specify the service with the
`--service` flag and the resource using the `--resource` flag:

```
$ aws-api-tool schema --service sqs --resource Queue openapi
properties:
  Attributes:
    additionalProperties: true
    type: object
  QueueName:
    type: string
  QueueUrl:
    type: string
  tags:
    additionalProperties: true
    type: object
type: object
```

```
$ aws-api-tool schema --service sns --resource Topic openapi
properties:
  Attributes:
    additionalProperties: true
    type: object
  Name:
    type: string
  Tags:
    items:
      properties:
        Key:
          type: string
        Value:
          type: string
      type: object
    type: array
  TopicArn:
    type: string
type: object
```

Note that different AWS service APIs will represent the same things
differently. An example of this is shown above where the AWS SQS Queue resource
uses the lowercase name "tags" to refer to a simple `map[string]string` whereas
the AWS SNS Topic resource uses the CamelCased name "Tags" and uses a list of
objects with a "Key" and "Value" property.

**NOTE**: By default, the `aws-api-tool schema openapi` command outputs the
OpenAPI3 Schema as YAML. You can output condensed JSON instead using the
`--output json` flag:


```
$ aws-api-tool schema --service sqs --resource Queue openapi --format json
{"properties":{"Attributes":{"additionalProperties":true,"type":"object"},"QueueName":{"type":"string"},"QueueUrl":{"type":"string"},"tags":{"additionalProperties":true,"type":"object"}},"type":"object"}
```

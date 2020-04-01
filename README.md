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

Use the `aws-api-tool list-apis` command to list AWS services:

```
$ aws-api-tool list-apis
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
$ aws-api-tool list-apis --filter ebs,sdb
+-------+-------------+----------------------------+
| ALIAS | API VERSION |         FULL NAME          |
+-------+-------------+----------------------------+
| ebs   | 2019-11-02  | Amazon Elastic Block Store |
| sdb   | 2009-04-15  | Amazon SimpleDB            |
+-------+-------------+----------------------------+
```

### Get summary information about an API

To get summary information about a particular AWS service API, use the
`aws-api-tool info <api>` command.

```
$ aws-api-tool info ec2
Full name:        Amazon Elastic Compute Cloud
API version:      2016-11-15
Protocol:         ec2
Total operations: 400
Total resources:  52
Total objects:    1991
Total scalars:    383
Total payloads:   738
Total exceptions: 0
Total lists:      405
```

### List API operations

To list operations for an AWS service API, use the `aws-api-tool
list-operations <api>` command:

```
$ aws-api-tool list-operations sns
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
$ aws-api-tool list-operations sns --prefix Update,Delete
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
```

```
$ aws-api-tool list-operations s3 --method GET --prefix ListO
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

Use the `aws-api-tool list-resources <api>` command to list these resource objects:

```
$ aws-api-tool list-resources sqs
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
$ aws-api-tool list-operations apigateway --prefix Create
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
$ aws-api-tool list-resources apigateway
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

Use the `aws-api-tool list-objects <api>` command to list an API's objects. You
can use the `--type` flag to filter objects by type ("scalar", "list",
"payload" and "exception"):

```
$ aws-api-tool list-objects sns --type scalar
+---------------------------+-------------+-----------+
|           NAME            | OBJECT TYPE | DATA TYPE |
+---------------------------+-------------+-----------+
| AmazonResourceName        | scalar      | string    |
| Binary                    | scalar      | blob      |
| MapStringToString         | scalar      | map       |
| MessageAttributeMap       | scalar      | map       |
| PhoneNumber               | scalar      | string    |
| String                    | scalar      | string    |
| SubscriptionAttributesMap | scalar      | map       |
| TagKey                    | scalar      | string    |
| TagValue                  | scalar      | string    |
| TopicAttributesMap        | scalar      | map       |
| account                   | scalar      | string    |
| action                    | scalar      | string    |
| attributeName             | scalar      | string    |
| attributeValue            | scalar      | string    |
| authenticateOnUnsubscribe | scalar      | string    |
| boolean                   | scalar      | boolean   |
| delegate                  | scalar      | string    |
| endpoint                  | scalar      | string    |
| label                     | scalar      | string    |
| message                   | scalar      | string    |
| messageId                 | scalar      | string    |
| messageStructure          | scalar      | string    |
| nextToken                 | scalar      | string    |
| protocol                  | scalar      | string    |
| string                    | scalar      | string    |
| subject                   | scalar      | string    |
| subscriptionARN           | scalar      | string    |
| token                     | scalar      | string    |
| topicARN                  | scalar      | string    |
| topicName                 | scalar      | string    |
+---------------------------+-------------+-----------+
```

#### Show OpenAPI3 Schema for API resource

Use the `aws-api-tool schema <api> <resource>` command to display the OpenAPI3
Schema for a specific resource in an AWS Service API.

```
$ aws-api-tool schema sqs Queue
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
$ aws-api-tool schema sns Topic
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
          maxLength: 128
          minLength: 1
          type: string
        Value:
          maxLength: 256
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

**NOTE**: By default, the `aws-api-tool schema <api> <resource>` command outputs the
OpenAPI3 Schema as YAML. You can output condensed JSON instead using the
`--output json` flag:


```
$ aws-api-tool schema sqs Queue --format json
{"properties":{"Attributes":{"additionalProperties":true,"type":"object"},"QueueName":{"type":"string"},"QueueUrl":{"type":"string"},"tags":{"additionalProperties":true,"type":"object"}},"type":"object"}
```

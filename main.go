package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

const appName = "short-url"

func nameGenerator(roleName string) string {
	return strings.Join([]string{appName, roleName}, "-")
}
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		account, err := aws.GetCallerIdentity(ctx, nil, nil)
		if err != nil {
			return err
		}

		region, err := aws.GetRegion(ctx, nil, nil)

		if err != nil {
			return err
		}

		// Create an IAM role.
		role, err := iam.NewRole(ctx, nameGenerator("role"), &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
					}]
			}`),
		})
		if err != err {
			return err
		}

		cloudwatch, err := cloudwatch.NewLogGroup(ctx, nameGenerator("log-group"), &cloudwatch.LogGroupArgs{
			RetentionInDays: pulumi.Int(3),
		})

		// Attach a policy to allow writing logs to CloudWatch
		logPolicy, err := iam.NewRolePolicy(ctx, nameGenerator("role-policy"), &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
				"Version": "2012-10-17",
			    "Statement": [{
                    "Effect": "Allow",
                    "Action": [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    "Resource": "arn:aws:logs:*:*:*"
                }]
		}`)})
		if err != nil {
			return err
		}

		// Set arguments for constructing the Lambda Function resource.
		args := &lambda.FunctionArgs{
			Handler: pulumi.String("handler"),
			Role:    role.Arn,
			Runtime: pulumi.String("go1.x"),
			Code:    pulumi.NewFileArchive("./tmp/handler.zip"),
		}

		// Create the lambda using the args.
		lfunc, err := lambda.NewFunction(
			ctx,
			nameGenerator("lambda"),
			args,
			pulumi.DependsOn([]pulumi.Resource{logPolicy, cloudwatch}), // Make sure the role policy is created first
		)

		// Create a New API Gateway
		gw, err := apigateway.NewRestApi(ctx, nameGenerator("api-gw"), &apigateway.RestApiArgs{
			Name:        pulumi.String(nameGenerator("gw")),
			Description: pulumi.String("A simple gateway"),
			Policy: pulumi.String(`{
			  "Version": "2012-10-17",
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Principal": {
					"Service": "lambda.amazonaws.com"
				  },
				  "Effect": "Allow",
				  "Sid": ""
				},
				{
				  "Action": "execute-api:Invoke",
				  "Resource": "*",
				  "Principal": "*",
				  "Effect": "Allow",
				  "Sid": ""
				}
			  ]
			}`)},
		)
		if err != nil {
			return err
		}

		// Add a resource to the API Gateway.
		apirsc, err := apigateway.NewResource(ctx, nameGenerator("api"), &apigateway.ResourceArgs{
			RestApi:  gw.ID(),
			PathPart: pulumi.String("{proxy+}"),
			ParentId: gw.RootResourceId,
		})
		if err != nil {
			return err
		}

		// Add a method to the API Gateway.
		_, err = apigateway.NewMethod(ctx, "AnyMethod", &apigateway.MethodArgs{
			HttpMethod:    pulumi.String("ANY"),
			Authorization: pulumi.String("NONE"),
			RestApi:       gw.ID(),
			ResourceId:    apirsc.ID(),
		})
		if err != nil {
			return err
		}

		// Add an integration to the API Gateway.
		// This makes communication between the API Gateway and the Lambda function work
		_, err = apigateway.NewIntegration(ctx, "LambdaIntegration", &apigateway.IntegrationArgs{
			HttpMethod:            pulumi.String("ANY"),
			IntegrationHttpMethod: pulumi.String("POST"),
			ResourceId:            apirsc.ID(),
			RestApi:               gw.ID(),
			Type:                  pulumi.String("AWS_PROXY"),
			Uri:                   lfunc.InvokeArn,
		})
		if err != nil {
			return err
		}

		// Add a resource based policy to the Lambda function.
		// This is the final step and allows AWS API Gateway to communicate with the AWS Lambda function
		permission, err := lambda.NewPermission(ctx, "APIPermission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  lfunc.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, gw.ID()),
		}, pulumi.DependsOn([]pulumi.Resource{apirsc}))
		if err != nil {
			return err
		}

		// create a new deployment
		_, err = apigateway.NewDeployment(ctx, nameGenerator("deployment"), &apigateway.DeploymentArgs{
			Description:      pulumi.String("Short URL Generate Service"),
			RestApi:          gw.ID(),
			StageDescription: pulumi.String("Production"),
			StageName:        pulumi.String("p"),
		}, pulumi.DependsOn([]pulumi.Resource{apirsc, lfunc, permission}))
		if err != nil {
			return err
		}

		ctx.Export("url", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/p/", gw.ID(), region.Name))
		return nil
	},
	)
}

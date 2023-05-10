from aws_cdk import (
    aws_autoscaling as autoscaling,
    aws_ec2 as ec2,
    aws_ecs as ecs,
    aws_iam as iam,
    aws_ecs_patterns as ecs_patterns,
    App, CfnOutput, Stack, aws_dynamodb
)
from constructs import Construct





class BonjourFargate(Stack):

    def __init__(self, scope: Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # create dynamo db
        demo_table = aws_dynamodb.Table(
            self, "demo_table",
            table_name="demo_table",
            partition_key=aws_dynamodb.Attribute(
                name="id",
                type=aws_dynamodb.AttributeType.STRING
            )
        )

        # Create VPC and Fargate Cluster
        # NOTE: Limit AZs to avoid reaching resource quotas
        vpc = ec2.Vpc(
            self, "GorfVpc",
            max_azs=2
        )

        cluster = ecs.Cluster(
            self, 'GorfEc2Cluster',
            vpc=vpc
        )


        role = iam.Role(self, "GorfRole",
            assumed_by=iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
            managed_policies=[
                iam.ManagedPolicy.from_aws_managed_policy_name("AmazonDynamoDBFullAccess")
            ]
        )

        policy = iam.Policy(self, "GorfPolicy",
            policy_name="GorfPolicy",
            statements=[
                iam.PolicyStatement(
                    effect=iam.Effect.ALLOW,
                    actions=["dynamodb:Query"],
                    resources=[f"{demo_table.table_arn}"]
                    # resources=[f"arn:aws:dynamodb:eu-north-1:198104462023:table/demo_table"]
                )
            ]
        )
        policy.attach_to_role(role)

        container_image = ecs.ContainerImage.from_asset("gorf_aws")
        fargate_service = ecs_patterns.ApplicationLoadBalancedFargateService(
            self, "FargateServiceForGorf",
            cluster=cluster,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                image=container_image,
                execution_role=role,
                task_role=role
            ),
            runtime_platform=ecs.RuntimePlatform(
                    operating_system_family=ecs.OperatingSystemFamily.LINUX,
                    cpu_architecture=ecs.CpuArchitecture.ARM64
                )
        )

        fargate_service.target_group.configure_health_check(
            path="/hello"
        )

        demo_table.grant_read_write_data(role)

        


        fargate_service.service.connections.security_groups[0].add_ingress_rule(
            peer = ec2.Peer.ipv4(vpc.vpc_cidr_block),
            connection = ec2.Port.tcp(80),
            description="Allow http inbound from VPC"
        )

        CfnOutput(
            self, "LoadBalancerDNSForGorf",
            value=fargate_service.load_balancer.load_balancer_dns_name
        )

app = App()
BonjourFargate(app, "Gorf")
app.synth()

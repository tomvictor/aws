from aws_cdk import (
    aws_autoscaling as autoscaling,
    aws_ec2 as ec2,
    aws_ecs as ecs,
    aws_ecs_patterns as ecs_patterns,
    App, CfnOutput, Stack
)
from constructs import Construct





class BonjourFargate(Stack):

    def __init__(self, scope: Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

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

        container_image = ecs.ContainerImage.from_asset("gorf_aws")
        fargate_service = ecs_patterns.ApplicationLoadBalancedFargateService(
            self, "FargateServiceForGorf",
            cluster=cluster,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                image=container_image
            ),
            runtime_platform=ecs.RuntimePlatform(
                    operating_system_family=ecs.OperatingSystemFamily.LINUX,
                    cpu_architecture=ecs.CpuArchitecture.ARM64
                )
        )

        fargate_service.target_group.configure_health_check(
            path="/hello"
        )


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

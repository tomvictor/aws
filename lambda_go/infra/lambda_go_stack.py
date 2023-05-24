from aws_cdk import (
    # Duration,
    Stack,
    aws_lambda as _lambda,
    aws_lambda_go_alpha,
)
from constructs import Construct


class LambdaGoStack(Stack):
    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        # The code that defines your stack goes here
        ENV = {}
        lambda_fun = aws_lambda_go_alpha.GoFunction(
            self,
            id="goapi",
            function_name="ApiCorsLambda",
            description="Demo Api Lambda",
            runtime=_lambda.Runtime.GO_1_X,
            entry="src/gorfapi",
            environment=ENV,
        )
        lambda_fun.add_function_url()

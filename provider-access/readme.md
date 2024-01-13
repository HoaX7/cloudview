## Secure authentication & useful links
- Steps taken to securely authenticate users with AWS, GCP, Azure.

## AWS
- Aws provides `cross account role` capabilities allowing us access to a clients aws accounts and running commands on their behalf.
- This allows us to securely connect to their accounts instead of them providing us access keys which is a bad practice.

# How are we connecting?
- Create a cloudformation stack template which allows us to create a role for cross account setup.
- Configure lambda function to handle callback from the stack to save user details such as account ID and role arn.
- We can use the role arn in the sdk to fetch resources.
- We need to upload the cloudformation template to AWS s3 bucket and must be in the same region as the stack creation.

# Useful links
- https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/crpg-ref-responses.html
- https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/crpg-ref-requesttypes-delete.html
- https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cfn-customresource.html#aws-resource-cloudformation-customresource-return-values

- Check the yaml/json file for cloud formation configurations.

# URL to quick create CF stack
<!-- https://us-east-1.console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/quickcreate?region=us-east-1&param_CloudFriendlyId=clr6ep9yu002kybpw4n6n4j0r&stackName=cloudfriendly-access-clr6ep9yu002kybpw4n6n4j0r&templateURL=https://cf-public-cf-templates.s3.amazonaws.com/cloudformation.json -->

- The cf stack creation might stall if the lambda function does not return data - so we need to store metadata about failed operations and run cfn-response send method manually by making a http request to the `responseUrl`.

# IAM creds
- We are using IAM User to connect to client AWS accounts to assumeRole to fetch temporary credentials.
```AWS access portal URL: https://d-9067f33598.awsapps.com/start, Username: kevin, One-time password: Ow4QWTu$<2O```
- arn:aws:iam::590183913396:user/kevin
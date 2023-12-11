
import { PUBLIC_GITHUB_OAUTH_CLIENT_ID } from "$env/static/public";

export const KONVA_KEYS = { STAGE_CTX: Symbol() };

export const OAUTH_URLS = {
	GITHUB: "https://github.com/login/oauth/" +
    `authorize?client_id=${PUBLIC_GITHUB_OAUTH_CLIENT_ID}&scope=user`
};

export const COOKIE_NAME = "cv-token";

export const DEFAULT_REGION = "us-west-2";


export const NAVBAR_HEIGHT = 64;

// Values must be same as backend app
export const AWS_SERVICES = {
	EC2: "ec2",
	APIGATEWAYV2: "apigatewayV2",
	ROUTE53: "route53",
	LAMBDA: "lambda",
	EKS: "eks",
	ELBV2: "elbV2",
	S3: "s3",
	EFS: "efs",
	RDS: "rds",
	CLOUDFRONT: "cloudfront",
	DYNAMODB: "dynamodb"
} as const;

export const TEXT_COLORS = {
	ec2: {
		bg: "#FF9900",
		color: "#FF9900"
	},
	lambda: {
		bg: "#FF9900",
		color: "#FF9900"
	},
	s3: {
		bg: "#569A31",
		color: "#569A31"
	},
	rds: {
		bg: "#004479",
		color: "#004479"
	},
	elbV2: {
		bg: "#7233FF",
		color: "#7233FF"
	}
};

export const LEGEND_NAMES = {
	VPC: "Vpc",
	SECURITY_GROUP: "Security Group",
	SUBNET: "Subnet"
};

export const NODE_POSITIONS = {
	LEFT: "left",
	RIGHT: "right",
	BOTTOM: "bottom",
	TOP: "top",
	OVERLAP: "overlap",
};

export const AWS_INTERNAL_SERVICES = [
	AWS_SERVICES.EC2,
	AWS_SERVICES.LAMBDA,
	AWS_SERVICES.S3,
	AWS_SERVICES.EFS,
	AWS_SERVICES.DYNAMODB,
	AWS_SERVICES.RDS,
	AWS_SERVICES.EKS,
];
export const AWS_EXTERNAL_SERVICES = [
	AWS_SERVICES.APIGATEWAYV2,
	AWS_SERVICES.CLOUDFRONT,
	AWS_SERVICES.ELBV2,
	AWS_SERVICES.ROUTE53,
];

export const STATUS_COLORS = {
	RUNNING: "#66ff66",
	STOPPED: "#ff0000",
	UNKNOWN: "#808080"
};
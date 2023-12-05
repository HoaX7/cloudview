import type { ConnectorTargetProps } from "./Konva";
import type { RDSProps } from "./aws/rds";

export type ServiceProps = {
  id: string;
  name: string;
  projectId: string;
  provider: string;
  accessKeySecret?: string;
  accessKeyId?: string;
  description?: string;
};

export type MetricDataReturnType = {
  name: string;
  result: AWSServices | any; // this prop contains values from ec2, lambda etc...
}[];

export type AWSServices =
  | ApiGatewayV2Props
  | Ec2Props
  | LambdaFunctionProps
  | S3Props
  | Route53Props
  | EKSProps
  | EFSProps
  | ELBV2Props
  | RDSProps;

// This type of properties is returned by aws sdk
export type ApiGatewayV2Props = {
  Items: {
    ApiId: string;
    ApiEndpoint: string;
    Name: string;
    ProtocolType: string;
  }[];
};

export type ApiGatewayV2IntegrationProps = {
  Items: {
    IntegrationUri: string;
    Description: string;
    IntegrationMethod: string;
    IntegrationId: string;
  }[];
};

export type ApiGatewayWithIntegrationProps = ApiGatewayV2Props["Items"][0] & {
  integrations: ApiGatewayV2IntegrationProps["Items"];
  lineTargets: ConnectorTargetProps[];
};

export type LambdaFunctionProps = {
  Functions: {
    FunctionName: string;
    Architectures: string[];
    FunctionArn: string; // unique id
    Handler: string;
    MemorySize: number;
    EphemeralStorage: {
      Size: number;
    };
    LastModified: Date;
    Runtime: string;
    State: string;
    StateReason: string | null;
    Role: string;
  }[];
  // has more props
};

export type S3Props = {
  Data: {
    Buckets: {
      CreationDate: Date;
      Name: string;
    }[];
  };
  Metrics: {
    Name: string;
    Statistics: {
      Datapoints: {
        Unit: string;
        Sum: number; // data in bytes
        Timestamp: Date;
      }[];
      Label: string;
    }
  }[];
  ACLList: {
    Bucket: string;
    ACL: {
      Grants: {
        Grantee: {
          Type: string;
          DisplayName: string;
          URI: string | null;
          ID: string;
          EmailAddress: string | null;
        };
        Permission: string;
      }[]
    }
  }[];
  // has more props
};

export type Ec2InstanceProps = {
  Architecture: string;
  BlockDeviceMappings: {
    DeviceName: string;
    Ebs: {
        VolumeId: string;
    }
  }[];
  CpuOptions: {
    AmdSevSnp: string;
    CoreCount: number;
    ThreadsPerCore: number;
  };
  InstanceType: string;
  Ipv6Address?: string;
  PublicIpAddress: string;
  PublicDnsName: string;
  PrivateIpAddress: string;
  PrivateDnsName: string;
  State: {
    Code: number;
    Name: string;
  };
  StateReason?: string;
  VpcId: string;
  KeyName: string;
  LaunchTime: string;
  NetworkInterfaces: any;
  Monitoring: { State: string };
  Placement: {
    AvailabilityZone: string;
  };
  SecurityGroups: {
    GroupId: string;
    GroupName: string;
  }[];
  SubnetId: string;
  ImageId: string; // use this to fetch os details
  InstanceId: string;
  Tags: {
    Key: string;
    Value: string;
  }[]
  // has more props
};

export type Ec2VolumeProps = {
    Size: number;
    VolumeType: string;
    State: string;
    Encrypted: string;
    VolumeId: string;
    Attachments: {
        Device: string;
    }[]
};
export type Ec2Props = {
  Reservations: {
    Instances: Ec2InstanceProps[];
  }[];
  Volumes: Ec2VolumeProps[];
};

export type ELBV2Props = {
  LoadBalancers: {
    AvailabilityZones: {
      LoadBalancerAddresses: [],
      SubnetId: string;
      ZoneName: string;
    }[];
    CanonicalHostedZoneId: string;
    CreatedTime: string;
    DNSName: string;
    IpAddressType: string;
    LoadBalancerArn: string;
    LoadBalancerName: string;
    Scheme: string;
    State: {
      Code: string;
      Reason: null | string;
    };
    VpcId: string;
    SecurityGroups: string[];
  }[];
};

export type EKSProps = {
  Clusters: [];
};

export type EFSProps = {
  FileSystems: {
    Name: string;
    CreationTime: string;
    CreationToken: string;
    FileSystemId: string;
    SizeInBytes: {
      Value: number;
    };
    FileSystemArn: string;
    PerformanceMode: string;
    Tags: {
      Key: string;
      Value: string;
    }[]
  }[];
};

export type Route53Props = {
  HostedZones: [];
};

type MetricProps = {
  MetricDataResults: {
    Id: string;
    Label: string;
    Timestamps: Date[];
    Values: number[];
  }[];
}
export type UsageProps = {
  Bandwidth: MetricProps;
  Cpu: MetricProps;
}

export type GroupedData = {
  externalGroup: MetricDataReturnType;
  internalGroup: MetricDataReturnType;
};
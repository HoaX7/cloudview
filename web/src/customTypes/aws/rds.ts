
export type RDSInstanceProps = {
    DBInstanceIdentifier: string;
    AllocatedStorage: number;
    AvailabilityZone: string;
    CACertificateIdentifier: string;
    CertificateDetails: {
        CAIdentifier: string;
        ValidTill: Date;
    };
    DBInstanceArn: string;
    DBInstanceClass: string;
    DBInstanceStatus: string;
    Endpoint: {
        "Address": string;
        "HostedZoneId": string;
        "Port": number;
    };
    Engine: string;
    EngineVersion: string;
    MasterUsername: string;
    PubliclyAccessible: boolean;
    StorageType: string;
    StorageThroughput: number;
    DBSubnetGroup: {
        Subnets: {
            SubnetAvailabilityZone: { Name: string; };
            SubnetIdentifier: string;
            SubnetStatus: string;
        }[];
        VpcId: string;
    };
    MaxAllocatedStorage: number;
    AutomaticRestartTime: Date;
    InstanceCreateTime: Date;
    VpcSecurityGroups: {
        VpcSecurityGroupId: string;
        Status: string;
    }[];
}

export type RDSProps = {
    DBInstances: RDSInstanceProps[];
}
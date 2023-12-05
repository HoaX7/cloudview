export type CloudFrontItemProps = {
    Id: string;
    ARN: string;
    DomainName: string;
    Enabled: boolean;
    Origins: {
        Items: {
            DomainName: string;
            Id: string;
            OriginShield: {
                OriginShieldRegion: string;
                Enabled: boolean;
            }
        }[];
    };
    Restrictions: {
        GeoRestriction: {
            RestrictionType: string;
            Items: null;
        };
    };
    Status: string;
    HttpVersion: string;

};

export type CloudFrontProps = {
    Items: CloudFrontItemProps[];
}
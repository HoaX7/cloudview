export type ProviderAccountProps = {
    id: string;
    name: string;
    projectId: string;
    provider: string;
    accessKeySecret?: string;
    accessKeyId?: string;
    description?: string;
    isDeleted?: boolean;
    metadata?: Record<string, unknown>;
    accountId: string;
    type: "ACCESS KEYS" | "CROSS ACCOUNT ROLE";
    accessRole?: string;
    featureAccessPermission: string;
  };

export type ProviderAccountWithProjectProps = {
    id: string;
    name: string;
    project: {
      id: string;
      name: string;
    };
    metadata?: Record<string, unknown>;
  }
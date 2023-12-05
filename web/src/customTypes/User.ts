export type UserProps = {
    id: string;
    username: string;
    email: string;
    avatarUrl: string;
    metadata?: Record<string, unknown>;
    createdAt: Date;
    updatedAt: Date;
    lastLoginAt: Date;
    isDeleted: boolean;
    subscriptionDaysLeft?: number;
    subscribedSince?: Date;
}

export type SessionProps = UserProps & {
    accessToken: string;
    provider: "github" | "google";
}

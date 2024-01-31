import type { UserProps } from "./user";

export type ProjectMemberProps = {
    id: string;
    projectId: string;
    userId: string;
    createdAt: Date;
    updatedAt: Date;
    isOwner: boolean;
    isActive: boolean;
    isDeleted: boolean;
}

export type ProjectMemberApiProps = ProjectMemberProps & { user: UserProps; };

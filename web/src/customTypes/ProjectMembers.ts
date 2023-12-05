import type { UserProps } from "./User";

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

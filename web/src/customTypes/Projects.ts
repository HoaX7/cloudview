import type { ServiceProps } from "./Services";

export type ProjectProps = {
    id: string;
    name: string;
    description?: string;
    email?: string;
    ownerId: string;
    members?: Record<string, unknown>;
    memberLimit: number;
    type: "PRIVATE" | "PUBLIC";
    metadata?: Record<string, unknown>;
    createdAt: Date;
    updatedAt: Date;
    ownerUsername?: string;
    isDeleted: boolean;
}

export type ProjectWithServiceProps = {
	project: ProjectProps;
	service: ServiceProps;
};

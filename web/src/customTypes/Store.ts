import type { ConnectableNodeProps, TargetFromNodeProps } from "./Konva";

export type DatastoreProps = {
    selectedRegion: string;
    konvaTargetFromNodes: TargetFromNodeProps[];
    konvaConnectableNodes: ConnectableNodeProps[];
    fetchData: boolean;
    dragNodeId: string | null;
}
import type { ConnectableNodeProps, TargetFromNodeProps } from "./Konva";

export type DatastoreProps = {
    selectedRegion: string;
    konvaTargetFromNodes: TargetFromNodeProps[];
    konvaConnectableNodes: ConnectableNodeProps[];
    fetchData: boolean;
    dragNodeId: string | null;
}

type Proportion = {
    x: number;
    y: number;
    width: number;
    height: number;
}
export type KonvaStoreProps = {
    externalBoundingRect: Proportion;
    internalBoundingRect: Proportion;
    rowCount: {
        internal: number,
        external: number
    }
}
export type SettingStoreProps = {
    animate: boolean;
}
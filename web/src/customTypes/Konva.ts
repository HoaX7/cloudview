import type { Stage } from "konva/lib/Stage";
import type { Vector2d } from "konva/lib/types";

export type StageContextProps = {
    getStage: () => Stage;
}

export type ConnectorTargetProps = {
    from: string;
    to: string;
}

export type TargetFromNodeProps = {
    id: string;
    x: number;
    y: number;
    label: string;
    lineStyle?: {
        dash?: number[];
        dashOffset?: number;
        stroke?: string;
    }
} & ConnectorTargetProps;

export type ConnectableNodeProps = Vector2d & { id: string; label: string; };

export type LegendProps = {
	id: string;
	count: number;
	name: string;
    highlight: string[];
    colors: string[];
  }

export type Points = Vector2d & { p?: 1 | 2 | 3 | 4; position?: string };

export type PathPoint = {
    command: string;
    x: string | number;
    y: string | number;
    position?: string;
  };

export type HighLightProps = {
    lines: string[];
    nodes: string[];
}
import type { PathPoint, Points } from "$src/customTypes/Konva";
import { NODE_POSITIONS } from "../constants";

const DEFAULT_CONTROL_POINTS = {
	dx: "",
	dy: "",
	x1: 0,
	y1: 0,
	x2: 0,
	y2: 0,
};
type CurveProps = {
  position: string;
  point1: number[];
  point2: number[];
};
export const calculateCurveControlPoints = (
	edge: PathPoint,
	nextEdge: PathPoint,
	position: string,
	prevPosiiton?: string
) => {
	const point1 = [ +edge.x, +edge.y ];
	const point2 = [ +nextEdge.x, +nextEdge.y ];
	/**
   * ##IMPORTANT INFORMATION TO UNDERSTAND THE LOGIC
   * 'prevPosition' holds the position of the previous edge or start point.
   * Ex: if the start point is below the edge the 'prevPosition' of the edge
   * will be 'top' therefore, the curve must be drawn from 'bottom'.
   * to draw the second part of the curve we must consider the position of
   * the current edge which is stored in 'position'. Ex: current edge is
   * on the left side from previous edge position. therefore, we must draw
   * the curve from bottom-left.
   * To understand better see a visualisation of the example:
   * https://codepen.io/hoax777/pen/WNPYMYZ
   */
	if (prevPosiiton === NODE_POSITIONS.TOP) {
		return drawTopCurve({
			position,
			point1,
			point2,
		});
	} else if (prevPosiiton === NODE_POSITIONS.LEFT) {
		return drawLeftCurve({
			position,
			point1,
			point2,
		});
	} else if (prevPosiiton === NODE_POSITIONS.RIGHT) {
		return drawRightCurve({
			point1,
			point2,
			position,
		});
	} else if (prevPosiiton === NODE_POSITIONS.BOTTOM) {
		return drawBottomCurve({
			point1,
			point2,
			position,
		});
	}
	return {
		...DEFAULT_CONTROL_POINTS,
		x1: point1[0],
		y1: point1[1],
		x2: point2[0],
		y2: point2[1],
		position,
	};
};
export const getEdgePosition = (edge1: PathPoint, edge2: PathPoint) => {
	const 
		dx = +edge2.x - +edge1.x;
	const dy = +edge2.y - +edge1.y;
	if (dx < 0) return NODE_POSITIONS.LEFT;
	else if (dx > 0) return NODE_POSITIONS.RIGHT;
	else if (dy < 0) return NODE_POSITIONS.TOP;
	else if (dy > 0) return NODE_POSITIONS.BOTTOM;
	return NODE_POSITIONS.OVERLAP;
};
export const calculateEdgePoints = (from: Points, to: Points) => {
	const axis = getControlAxis(from.p || 1);
	const midPoint = (from[axis] + to[axis]) / 2;
	const useMidPoint = true;
	const otherAxis = axis === "x" ? "y" : "x";
	return {
		command: "L",
		y: 0,
		x: 0,
		[axis]: useMidPoint ? midPoint : from[axis],
		[otherAxis]: from[otherAxis],
	};
};

const drawBottomCurve = (params: CurveProps) => {
	const { position, point1, point2 } = params;
	// eslint-disable-next-line prefer-const
	let [ x1, y1 ] = point1;
	// eslint-disable-next-line prefer-const
	let [ x2, y2 ] = point2;
	let dx = "";
	let dy = "";
	if (position === NODE_POSITIONS.RIGHT) {
		const path = edgeCurvePaths.topToRight();
		dx = path.dx;
		dy = path.dy;
		y1 = y1 + path.y[0];
	} else if (position === NODE_POSITIONS.LEFT) {
		const path = edgeCurvePaths.topToLeft();
		dx = path.dx;
		dy = path.dy;
		y1 = y1 + path.y[0];
	}
	return {
		x1,
		x2,
		y1,
		y2,
		dx,
		dy,
		position,
	};
};
const drawTopCurve = (params: CurveProps) => {
	const { position, point1, point2 } = params;
	// eslint-disable-next-line prefer-const
	let [ x1, y1 ] = point1;
	// eslint-disable-next-line prefer-const
	let [ x2, y2 ] = point2;
	let dx = "";
	let dy = "";
	if (position === NODE_POSITIONS.LEFT) {
		const path = edgeCurvePaths.bottomToLeft();
		y1 = y1 + path.y[0];
		dx = path.dx;
		dy = path.dy;
	} else if (position === NODE_POSITIONS.RIGHT) {
		const path = edgeCurvePaths.bottomToRight();
		dx = path.dx;
		dy = path.dy;
		y1 = y1 + path.y[0];
	}
	return {
		x1,
		y1,
		x2,
		y2,
		dx,
		dy,
	};
};
const drawRightCurve = (params: CurveProps) => {
	const { position, point1, point2 } = params;
	// eslint-disable-next-line prefer-const
	let [ x1, y1 ] = point1;
	// eslint-disable-next-line prefer-const
	let [ x2, y2 ] = point2;
	let dx = "";
	let dy = "";
	if (position === NODE_POSITIONS.TOP) {
		const path = edgeCurvePaths.leftToTop();
		x1 = x1 + path.x[0];
		dx = path.dx;
		dy = path.dy;
	} else if (position === NODE_POSITIONS.BOTTOM) {
		const path = edgeCurvePaths.leftToBottom();
		dx = path.dx;
		dy = path.dy;
		x1 = x1 + path.x[0];
	}
	return {
		x1,
		y1,
		x2,
		y2,
		dx,
		dy,
	};
};
const drawLeftCurve = (params: CurveProps) => {
	const { position, point1, point2 } = params;
	// eslint-disable-next-line prefer-const
	let [ x1, y1 ] = point1;
	const [ x2, y2 ] = point2;
	let dx = "";
	let dy = "";
	if (position === NODE_POSITIONS.TOP) {
		const path = edgeCurvePaths.rightToTop();
		x1 = x1 + path.x[0];
		dx = path.dx;
		dy = path.dy;
	} else if (position === NODE_POSITIONS.BOTTOM) {
		const path = edgeCurvePaths.rightToBottom();
		dx = path.dx;
		dy = path.dy;
		x1 = x1 + path.x[0];
	}
	return {
		x1,
		y1,
		x2,
		y2,
		dx,
		dy,
	};
};
// cp1 & cp2 values (control points)
const edgeCurvePaths = {
	// values for smoother curve - 10 for dx, dy and 20 for differences.
	bottomToRight: () => ({
		dx: "0 -8",
		dy: "8 -8",
		x: [ 0, 0 ],
		y: [ 16, 0 ],
	}),
	bottomToLeft: () => ({
		dx: "0 -8",
		dy: "-8 -8",
		x: [ 0, 0 ],
		y: [ 16, 0 ],
	}),
	topToRight: () => ({
		dx: "0 8",
		dy: "8 8",
		x: [ 0, 0 ],
		y: [ -16, 0 ],
	}),
	topToLeft: () => ({
		dx: "0 8",
		dy: "-8 8",
		x: [ 0, 0 ],
		y: [ -16, 0 ]
	}),
	leftToBottom: () => ({
		dx: "8 0",
		dy: "8 8",
		x: [ -16, 0 ],
		y: [ 0, 0 ],
	}),
	leftToTop: () => ({
		dx: "8 0",
		dy: "8 -8",
		x: [ -16, 0 ],
		y: [ 0, 0 ],
	}),
	rightToTop: () => ({
		dx: "-8 0",
		dy: "-8 -8",
		x: [ 16, 0 ],
		y: [ 0, 0 ],
	}),
	rightToBottom: () => ({
		dx: "-8 0",
		dy: "-8 8",
		x: [ 16, 0 ],
		y: [ 0, 0 ]
	})
};

const getControlAxis = (position: number) => {
	if (position % 2 === 0) return "x";
	return "y";
};
import type { Vector2d } from "konva/lib/types";
import type { DatastoreProps } from "$src/customTypes/Store";
import { NODE_POSITIONS } from "../constants";
import type { PathPoint, Points } from "$src/customTypes/Konva";
import { calculateCurveControlPoints, calculateEdgePoints, getEdgePosition } from "./edgePaths";
import { get } from "svelte/store";
import KonvaStore from "$src/store/konva";

// This function is used to fetch line points to
// draw arrow between services
/**
 * This function returns konva 'Line' points without edges and curves.
 * @param from
 * @param to
 * @returns
 */
export const getConnectorPoints = (from: Vector2d, to: Vector2d) => {
	const dx = to.x - from.x;
	const dy = to.y - from.y;
	const angle = Math.atan2(-dy, dx);

	// Increase radius to add a gap between
	// icon and leading/tailing arrow
	const radius = 0;

	return [
		from.x + -radius * Math.cos(angle + Math.PI),
		from.y + radius * Math.sin(angle + Math.PI),
		to.x + -radius * Math.cos(angle),
		to.y + radius * Math.sin(angle),
	];
};

export const getArrowHeadPoints = (endPoint: Points) => {
	if (endPoint.p === 1) {
		return getConnectorPoints({
			x: endPoint.x,
			y: endPoint.y - 20
		}, endPoint);
	} else if (endPoint.p === 2) {
		return getConnectorPoints({
			x: endPoint.x + 20,
			y: endPoint.y
		}, endPoint);
	} else if (endPoint.p === 3) {
		return getConnectorPoints({
			x: endPoint.x,
			y: endPoint.y + 20
		}, endPoint);
	} else if (endPoint.p === 4) {
		return getConnectorPoints({
			x: endPoint.x - 20,
			y: endPoint.y
		}, endPoint);
	}
	return [];
};

export const stripIntegrationUriIp = (str: string) => {
	// Use a regular expression to extract the IP address
	const ipRegex = /\/\/([^:/]+):/;
	const match = str.match(ipRegex);

	// Extracted IP address is in match[1]
	return match ? match[1] : str;
};

/**
 * Responsible to calculate the start-end positions at which the arrows must
 * be drawn.
 * Each arrow is drawn from the center of icon on each edge.
 * We are using the width and height of the icon as reference points to traverse through each edge
 * and X,Y axis to find out the placement of the target.
 *
 * "p" indicates the 4 points on the edges of the rectangle used to connect
 * the line dots.
 * p1 top, p2 right, p3 bottom, p4 left. the point moves in a clockwise direction.
 * @param from
 * @param to
 * @returns
 */

export const getConnectorPointsByPosition = (from: Points, to: Points) => {
	const width = 80; // width of the icon
	const height = 80; // height of the icon

	const dx = to.x - from.x;
	const dy = to.y - from.y;

	const position = getTargetPosition(dx, dy);
	adjustHorizontal(position, from, to, width, height);
	adjustVertical(position, from, to, width, height);

	// Only setting target position to draw edge curves.
	to.position = position;

	return {
		from,
		to,
	};
};

/**
 * This function returns the connection mappings used to draw lines
 * from different services.
 * We can use this data to show a mapping table to show
 * what services are connecting to/from the selected service.
 * @param datastore
 * @param id
 * @returns
 */
export const getConnectorMappings = (
	datastore: DatastoreProps,
	id?: string
) => {
	const connections: {
    from: string;
    to: string;
  }[] = [];
	if (!id) return connections;
	(datastore?.konvaTargetFromNodes || []).forEach((tg) => {
		if (tg.from === id) {
			const connectingTo = datastore.konvaConnectableNodes.filter(
				(node) => node.id === tg.to || tg.to.includes(node.id)
			);
			connectingTo.forEach((c) => {
				connections.push({
					from: tg.label || tg.id,
					to: c.label || c.id,
				});
			});
		}
		if (id && tg.to.includes(id)) {
			const node = datastore.konvaConnectableNodes.find((nd) => nd.id === id);
			if (node) {
				connections.push({
					from: tg.label || tg.id,
					to: node.label || node.id,
				});
			}
		}
	});

	return connections;
};

/**
 * To place external services in column layout the `offset` property is used
 * to get the number of items that are above the item that is currently being placed.
 * Example 1: 2 api services are placed 1 below the other, now to place cloudfront
 * below these 2 services we need an offset of 2.
 * Example 2: If there are 2 api services and 1 cloudfront services to place
 * Load balancer service the offset is 3.
 * `offset` is the count of all previous items before the current item.
 * @param offset
 * @param itemIndex
 * @param group
 * @returns
 */
export const getProportions = (
	offset: number,
	itemIndex: number,
	group: "external" | "internal",
	defaultX = 120,
	defaultY = 120
) => {
	const konvastore = get(KonvaStore.getStore());
	let x = group === "internal" ? konvastore.internalBoundingRect.x : konvastore.externalBoundingRect.x;
	let y = group === "external" ? konvastore.externalBoundingRect.y : defaultY;
	const imageWidth = 100; // This includes text above image. (+20)
	const imageHeight = 100;
	const xPadding = 10;
	const yPadding = 80;
	if (group === "internal") {
		x = x + (imageWidth + xPadding) * itemIndex;
		y = y + (imageHeight + yPadding) * offset;
	} else if (group === "external") {
		y =
      y +
      (imageHeight + xPadding) * itemIndex +
      (imageHeight + yPadding) * offset;
	}

	return {
		x,
		y,
	};
};


export const drawSVGPath = (from: Points, to: Points) => {
	const points: PathPoint[] = [];
	const startPoint: PathPoint = {
		command: "M",
		x: from.x,
		y: from.y,
	};
	const endPoint: PathPoint = {
		command: "L",
		x: to.x,
		y: to.y,
	};
	if (from.x === to.x || from.y === to.y) {
		points.push(startPoint, endPoint);
	} else {
		// There is always either going to be 2 or 3 edges.
		const edge1 = calculateEdgePoints(from, to);
		const edge2 = calculateEdgePoints(to, from);
		const extraEdges: PathPoint[] = [];
		const edges: PathPoint[] = [ startPoint, edge1, edge2, ...extraEdges, endPoint ];
		const edgesWithCurves: PathPoint[] = [];
		edges.forEach((edge, index) => {
			const nextIndex = index + 1;
			const nextEdge = edges[nextIndex];
			if (!nextEdge) {
				edgesWithCurves.push(edge);
				return;
			}
			/**
       * Use the control points to adjust the extra line points
       * on `x` and `y` axis. This ensures the curve is drawn
       * correctly over the lines.
       *
       * Reason for re-adjusting x and y axis is to keep the line straight
       * at all times.
       */

			// This position will help us determine the next
			// edge position to draw the curve from
			// bottom-right or top-left, top-right, bottom-left etc...
			// let nextPosition: string = "";
			// if (edges[nextIndex + 1]) {
			// 	nextPosition = getEdgePosition(nextEdge, edges[nextIndex + 1]);
			// }

			// 'position' implies the placement of control points along y or x axis.
			// dx != 0 -> x axis. dy != 0 -> y axis
			const position = getEdgePosition(edge, nextEdge);
			edge.position = position;
			edges[index].position = position;
			if (index === 0) {
				edgesWithCurves.push(edge);
				return;
			}

			const controlPoints = calculateCurveControlPoints(
				edge,
				nextEdge,
				position,
				edges[index - 1]?.position
			);

			edge.x = controlPoints.x1;
			edge.y = controlPoints.y1;
			edges[nextIndex].x = controlPoints.x2;
			edges[nextIndex].y = controlPoints.y2;

			edgesWithCurves.push(edge);

			if (controlPoints.dx !== "" && controlPoints.dy !== "") {
				edgesWithCurves.push({
					command: "t",
					x: controlPoints.dx,
					y: controlPoints.dy
				});
			}
		});
		points.push(...edgesWithCurves);
	}

	return points.map((p) => `${p.command}${p.x} ${p.y}`.trim()).join(" ");
};

const adjustHorizontal = (
	position: string,
	from: Points,
	to: Points,
	w: number,
	h: number
) => {
	if (position === NODE_POSITIONS.RIGHT) {
		adjustPoint(from, w, h / 2);
		from.p = 2;
		to.y = to.y + h / 2;
		to.p = 4;
	} else if (position === NODE_POSITIONS.LEFT) {
		from.y = from.y + h / 2;
		from.p = 4;
		adjustPoint(to, w, h / 2);
		to.p = 2;
	}
};
const adjustVertical = (
	position: string,
	from: Points,
	to: Points,
	w: number,
	h: number
) => {
	if (position === NODE_POSITIONS.TOP) {
		from.x = from.x + w / 2;
		from.p = 1;
		adjustPoint(to, w / 2, h);
		to.p = 3;
	} else if (position === NODE_POSITIONS.BOTTOM) {
		adjustPoint(from, w / 2, h);
		from.p = 3;
		to.x = to.x + w / 2;
		to.p = 1;
	}
};
const getTargetPosition = (dx: number, dy: number) => {
	// hard-coded icon dimensions
	const dimension = 80;
	if (dx + dimension < 0) return NODE_POSITIONS.LEFT;
	else if (dx > dimension) return NODE_POSITIONS.RIGHT;
	else if (dy + dimension < 0) return NODE_POSITIONS.TOP;
	else if (dy > dimension) return NODE_POSITIONS.BOTTOM;
	return NODE_POSITIONS.OVERLAP;
};
const adjustPoint = (point: Vector2d, offsetX: number, offsetY: number) => {
	point.x += offsetX;
	point.y += offsetY;
};
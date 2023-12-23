import { Vector2d } from "konva/lib/types";
import { COLOR_SCHEME, NODE_POSITIONS } from "./constants";
import Konva from "konva";
import { clone } from ".";
import { Container } from "konva/lib/Container";
import { Node, NodeConfig } from "konva/lib/Node";

type LineConnection = {
	from: Point;
	to: Point;
	p: number;
	position?: string;
	opacity: number;
	name: string;
}
export const drawKonvaCanvas = (el: HTMLDivElement, proportions: { width: number; height: number; }, options = {
	scale: .6,
	animate: false
}) => {
	const stage = new Konva.Stage({
		container: el,
		draggable: true,
		width: proportions.width,
		height: proportions.height,
		scale: {
			x: options.scale,
			y: options.scale
		}
	});
	const layer = new Konva.Layer();
	const clearLines = (idsToDelete: string[]) => {
		const children = layer
			?.getChildren()
			.filter((child) => idsToDelete.includes(child.attrs.id));
		children?.forEach((child) => child.destroy());
	};
	const drawLines = (linesArray: LineConnection[], animate = false) => {
		linesArray.map((line) => {
			if (!line) return;
			const arrowHeadPoints = getArrowHeadPoints({
				x: line?.to.x || 0,
				y: line?.to.y || 0,
				p: line?.p || 0
			});
			const path = drawSvgPath(line.from, line.to);
			const pathClone = svgPath.clone({
				data: path,
				id: `${line.from.x}-${line.to.x}`,
				opacity: line.opacity,
			});
			const circleClone = circle.clone({
				x: line?.from.x,
				y: line?.from.y,
				id: `${line.from.x}-${line.to.x}`,
				opacity: line.opacity,
			});
			const arrowClone = arrow.clone({
				points: arrowHeadPoints,
				id: `${line.from.x}-${line.to.x}`,
				opacity: line.opacity,
			});
			if (animate) {
				const pathLen = pathClone.getLength();
				pathClone.dash([ 5, 5 ]);
				pathClone.dashOffset(10);
				const anim = new Konva.Animation(function (frame) {
					const dashLen = pathLen - (frame?.time || 0) / 30;
					pathClone.dashOffset(dashLen);
					if (dashLen < 0) {
						// anim.stop()
						arrowClone.visible(true);
					}
				}, layer);
				anim.start();
			}
			layer.add(pathClone);
			layer.add(circleClone);
			layer.add(arrowClone);
		});
		layer.batchDraw();
	};
	const highlightLines = (lines: LineConnection[], attrs: { id: string; }) => {
		clearLines(lines.map((ln) => `${ln.from.x}-${ln.to.x}`));
		if ([ "Ec2", "Lambda" ].includes(attrs.id)) {
			attrs.id = "Api Gateway";
		}
		const selectedLines = clone(lines).map((l: LineConnection) => {
			if (l.name === attrs.id) {
				l.opacity = 1;
			} else {
				l.opacity = .3;
			}
			return l;
		});
		drawLines(selectedLines, false);
	};
	const highlightInstance = (id: string) => {
		const idToCompare = id.toLowerCase();
		if ([ "ec2", "lambda", "api gateway" ].includes(idToCompare)) {
			// highlight apigateway, ec2 and lambda
			const shape = stage.findOne("#Ec 2");
			shape?.opacity(1);
			const shape2 = stage.findOne("#Lambda");
			shape2?.opacity(1);
			const shape3 = stage.findOne("#CDN");
			shape3?.opacity(.4);
		} else {
			const shape = stage.findOne("#Ec 2");
			shape?.opacity(.4);
			const shape2 = stage.findOne("#Lambda");
			shape2?.opacity(.4);
			const shape3 = stage.findOne("#CDN");
			shape3?.opacity(1);
			const shape4 = stage.findOne("#Api Gateway");
			shape4?.opacity(1);
		}
	};
	const resetInstance = () => {
		const shape = stage.findOne("#Ec 2");
		shape?.opacity(1);
		const shape2 = stage.findOne("#Lambda");
		shape2?.opacity(1);
		const shape3 = stage.findOne("#CDN");
		shape3?.opacity(1);
		const shape4 = stage.findOne("#Api Gateway");
		shape4?.opacity(1);
	};

	const imageRect = new Konva.Rect({
		width: 100,
		height: 100,
		draggable: false,
		fill: "black",
		cornerRadius: 5,
	});
	const showPopup = (parent: Container<Node<NodeConfig>>, attrs: any) => {
		const rect = imageRect.clone({
			fill: attrs.fill,
			cornerRadius: 5,
			draggable: false,
			listening: false,
			width: 200,
			height: 60,
		});
		const popupGroup = new Konva.Group({
			listening: false,
			x: 0,
			y: -80,
			id: "popup"
		});
		const textContent = [];
		if (attrs.id === "Api Gateway") {
			textContent.push("ANY https:///12.22.34.400:8000",
				"POST arn:function:test-lambda");
		} else if (attrs.id === "Ec2") {
			textContent.push("Public ip: 12.22.34.400", "State: running");
		} else if (attrs.id === "Lambda") {
			textContent.push("runtime: nodejsx18.0", "memory size: 512MB");
		} else if (attrs.id === "CDN") {
			textContent.push("Http version: HTTP2", "Status: Deployed");
			popupGroup.x(-120);
			rect.width(150);
		}
		popupGroup.add(rect);
		textContent.map((text, index) => {
			const textEl = new Konva.Text({
				fill: "white",
				fontStyle: "bold",
				listening: false,
				text: text,
				x: 10,
				y: 15 + (index * 20)
			});
			popupGroup.add(textEl);
		});
		parent.add(popupGroup);
	};
	const resetPopup = (parent: Container<Node<NodeConfig>>) => {
		const child = parent.findOne("#popup");
		child?.destroy();
	};
    
	const svgPath = new Konva.Path({
		listening: false,
		perfectDrawEnabled: false,
		stroke: COLOR_SCHEME.CONNECTOR,
		strokeWidth: 3,
		draggable: false,
		data: "M0 0 L10 10",
		zIndex: 0,
	});
	const circle = new Konva.Circle({
		listening: false,
		draggable: false,
		radius: 6,
		stroke: COLOR_SCHEME.CONNECTOR,
		fill: "white",
		zIndex: 999,
		strokeWidth: 2,
		perfectDrawEnabled: false,
	});
	const arrow = new Konva.Arrow({
		draggable: false,
		listening: false,
		perfectDrawEnabled: false,
		points: [],
		stroke: COLOR_SCHEME.CONNECTOR,
		strokeWidth: 2,
		fill: COLOR_SCHEME.CONNECTOR,
	});
	const mainGroup = new Konva.Group({
		draggable: false,
		// listening: false,
		width: 1,
		height: 1,
	});
	const mainText = new Konva.Text({
		listening: false,
		fontStyle: "bold",
	});
	arrow.cache();
	circle.cache();
	svgPath.cache();
	const instances = getInstances(proportions);
	const linesToDraw = instances.map((item) => {
		const group = mainGroup.clone({
			id: item.name,
			...item.config
		});
		const text = mainText.clone({
			y: -14,
			text: item.name,
		});
		group.add(text);
		const clonedRect = imageRect.clone({
			fill: item.fill,
			id: item.name 
		});
		group.add(clonedRect);
		const imgObj = new Image();
		imgObj.onload = function() {
			const image = new Konva.Image({
				image: imgObj,
				listening: false,
				perfectDrawEnabled: false,
				x: 22,
				y: 22
			});
			group.add(image);
		};
		imgObj.src = item.icon;
		layer.add(group);
		group.on("mouseover", function (evt) {
			const shape = evt.target;
			const parentGroup = shape.getParent();
			document.body.style.cursor = "pointer";
			const attr = shape.getAttrs();
			highlightLines(linesToDraw as LineConnection[], attr);
			highlightInstance(attr.id);
			if (parentGroup) {
				showPopup(parentGroup, shape.attrs);
			}
		});
		group.on("mouseout", function (evt) {
			const shape = evt.target;
			document.body.style.cursor = "default";
			// reset highlight
			clearLines(linesToDraw.map((l) => `${l?.from.x}-${l?.to.x}`));
			drawLines(linesToDraw as LineConnection[], true);
			resetInstance();
			const parent = shape.getParent();
			if (parent) resetPopup(parent);
		});
		return item.connections;
	}).filter(Boolean).flat();

	console.log("drawing connector lines:", linesToDraw.length);
	drawLines(linesToDraw as LineConnection[], true);
	
	layer.batchDraw();
	stage.add(layer);
};

export const getInstances = (containerProportions: {
    width: number;
    height: number;
}) => {
	const midX = containerProportions.width / 2;
	const totalW = 500;
	const leftPadding = 100;
	const x = midX - (totalW / 2) + leftPadding;
	const internalServicesX = x + (totalW - 100);
	return [ {
		name: "Api Gateway",
		config: {
			x: x,
			y: 100
		},
		fill: COLOR_SCHEME.GATEWAY,
		icon: "/assets/images/aws/api-gateway.png",
		connections: [ {
			from: {
				x: x + 100,
				y: 150,
				p: 2
			},
			to: {
				x: internalServicesX,
				y: 150,
				p: 4
			},
			p: 4,
			opacity: 1,
			name: "Api Gateway",
		}, {
			from: {
				x: x + 100,
				y: 150,
				p: 2
			},
			to: {
				x: internalServicesX,
				y: 300,
				p: 4
			},
			p: 4,
			opacity: 1,
			name: "Api Gateway",
		} ]
	}, {
		name: "Ec2",
		config: {
			x: internalServicesX,
			y: 100
		},
		fill: COLOR_SCHEME.VM,
		icon: "/assets/images/aws/ec2.png",
	}, {
		name: "Lambda",
		config: {
			x: internalServicesX,
			y: 250
		},
		fill: COLOR_SCHEME.SERVERLESS,
		icon: "/assets/images/aws/lambda.png",
	}, {
		name: "CDN",
		config: {
			x: x,
			y: 300
		},
		fill: COLOR_SCHEME.CDN,
		icon: "/assets/images/aws/cloudfront.png",
		connections: [ {
			from: {
				x: x + 50,
				y: 300,
				p: 1
			},
			to: {
				x: x + 50,
				y: 200,
				p: 3
			},
			p: 3,
			opacity: 1,
			name: "CDN",
		} ]
	} ];
};

export const loadCanvas = (id: string, options = {
	scale: .6,
	animate: false
}) => {
	const el = document.getElementById(id) as HTMLDivElement;
	const parentContainer = el?.parentElement;
	console.log(parentContainer, parentContainer?.getBoundingClientRect());
	if (!parentContainer) {
		console.error("Parent element not found. InfrastructureCanvas must be wrapped inside a 'div' element");
	} else {
		drawKonvaCanvas(el, parentContainer.getBoundingClientRect(), options);
	}
};
type PathPoint = {
    command: string;
    x: number;
    y: number;
    posiiton?: string;
}
type Point = Vector2d & {
    p?: number;
    position?: string;
}
export const drawSvgPath = (from: Point, to: Point) => {
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
		const dx = to.x - from.x;
		const dy = to.y - from.y;
		to.position = getTargetPosition(dx, dy);
		const edge1 = calculateEdgePoints(from, to);
		const edge2 = calculateEdgePoints(to, from);
		points.push(startPoint, edge1, edge2, endPoint);
	}

	return points.map((p) => `${p.command}${p.x} ${p.y}`.trim()).join(" ");
};

export const calculateEdgePoints = (from: any, to: any) => {
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
export const getArrowHeadPoints = (endPoint: Vector2d & {
    p: number;
}) => {
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

const getConnectorPoints = (from: Vector2d, to: Vector2d) => {
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
const getControlAxis = (position: number) => {
	if (position % 2 === 0) return "x";
	return "y";
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
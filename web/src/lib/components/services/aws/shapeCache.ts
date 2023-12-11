import Konva from "konva";
import type { RectConfig } from "konva/lib/shapes/Rect";

const imageWidth = 80;
const imageHeight = 80;
const rect = new Konva.Rect({
	width: imageWidth,
	height: imageHeight,
	cornerRadius: 5,
	x: 0,
	y: 0,
	shadowBlur: 10,
	shadowOffset: {
		x: 5,
		y: 5
	}
});
rect.cache();

const boundingGroupRect = new Konva.Rect({
	width: 100,
	height: 145,
	cornerRadius: 5,
	x: 0,
	y: 0,
	listening: true,
});
boundingGroupRect.cache();

// We are caching the rect as the dimensions are static
export const getImageRect = (obj: RectConfig) => {
	obj.shadowColor = obj.fill;
	return rect.clone(obj);
};

export const getBoundingRect = (obj: RectConfig) => {
	return boundingGroupRect.clone(obj);
};

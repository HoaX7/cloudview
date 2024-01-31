import type { ResourceDataReturnType } from "$src/customTypes/services";
import KonvaStore from "$src/store/konva";
import {
	AWS_EXTERNAL_SERVICES,
	AWS_INTERNAL_SERVICES,
	AWS_SERVICES,
} from "./constants";
import { getProportions } from "./konva";

export const precomputeBorder = (data: ResourceDataReturnType) => {
	const konvastore = KonvaStore;
	const internalBoundingRect = computeDimensions(
		data.filter((d: any) => AWS_INTERNAL_SERVICES.includes(d.name)),
		"internal"
	);
	const externalBoundingRect = computeDimensions(
		data.filter((d: any) => AWS_EXTERNAL_SERVICES.includes(d.name)),
		"external"
	);
	const userConnectionWidth = 128;
	const spacing = 200; // spacing between external - internal services;
	const totalWidth =
    internalBoundingRect.width +
    spacing +
    externalBoundingRect.width +
    spacing +
    userConnectionWidth;

	const midX = window.innerWidth / 2;
	const x = midX - (totalWidth / 2); // center align along x axis
	externalBoundingRect.x = x + (userConnectionWidth + spacing);
	internalBoundingRect.x = x + (totalWidth - internalBoundingRect.width);
	externalBoundingRect.y = (internalBoundingRect.height / 2) - (externalBoundingRect.height / 2);
	if (internalBoundingRect.height === externalBoundingRect.height) {
		externalBoundingRect.y = internalBoundingRect.y;
	}
	konvastore.update({
		internalBoundingRect,
		externalBoundingRect,
		rowCount: {
			internal: internalBoundingRect.rowCount.internal,
			external: externalBoundingRect.rowCount.external
		}
	});
};
const computeDimensions = (
	data: ResourceDataReturnType,
	type: "internal" | "external"
) => {
	const imageWidth = 100;
	const imageHeight = 100;
	const imagePadding = 10;
	const imageYPadding = 85;
	const rowCount = {
		internal: 0,
		external: 0,
	};
	const numberOfInstances = data.map((item) => {
		if (type === "external") {
			if (item.name === AWS_SERVICES.APIGATEWAYV2) {
				rowCount.external += 1;
				return item.result.length;
			} else if (item.name === AWS_SERVICES.CLOUDFRONT) {
				rowCount.external += 1;
				return item.result.Items.length;
			} else if (item.name === AWS_SERVICES.ELBV2) {
				rowCount.external += 1;
				return item.result.LoadBalancers.length;
			} else if (item.name === AWS_SERVICES.ROUTE53) {
				rowCount.external += 1;
				return item.result.HostedZones.length;
			}
		}
		if (item.name === AWS_SERVICES.EC2) {
			rowCount.internal += 1;
			return item.result.Reservations.length;
		} else if (item.name === AWS_SERVICES.DYNAMODB) {
			rowCount.internal += 1;
			return item.result.length;
		} else if (item.name === AWS_SERVICES.EFS) {
			rowCount.internal += 1;
			return item.result.FileSystems.length;
		} else if (item.name === AWS_SERVICES.EKS) {
			rowCount.internal += 1;
			return item.result.Clusters.length;
		} else if (item.name === AWS_SERVICES.LAMBDA) {
			rowCount.internal += 1;
			return item.result.Functions.length;
		} else if (item.name === AWS_SERVICES.RDS) {
			rowCount.internal += 1;
			return item.result.DBInstances.length;
		} else if (item.name === AWS_SERVICES.ROUTE53) {
			rowCount.internal += 1;
			return item.result.HostedZones.length;
		} else if (item.name === AWS_SERVICES.S3) {
			rowCount.internal += 1;
			return item.result.Data.Buckets.length;
		}
		return 0;
	});

	let maxInstances = Math.max(...numberOfInstances);
	if (maxInstances < 2) maxInstances = 2;
	if (type === "internal") {
		const proportions = getProportions(0, 0, "internal");
		const width = maxInstances * (imageWidth + imagePadding);
		const height =
      rowCount.internal * imageHeight + rowCount.internal * imageYPadding;
		return {
			...proportions,
			width,
			height,
			rowCount
		};
	}
	const totalInstances = numberOfInstances.reduce((acc, r) => acc + r, 0);
	const proportions = getProportions(0, 0, "external");
	const width = imageWidth + imagePadding + (imageWidth / 2); // extra padding;
	const height = totalInstances * imageHeight + totalInstances * imageYPadding;
	return {
		...proportions,
		width,
		height,
		rowCount
	};
};

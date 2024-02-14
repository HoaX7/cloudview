<!-- @component
    This component renders aws service icons & its info.

    Component must be placed inside a `InfiniteCanvas` Component.
-->

<script lang="ts">
  import type {
  	ApiGatewayV2Props,
  	ApiGatewayWithIntegrationProps,
  	GroupedData,
  	ResourceDataReturnType,
  } from "$src/customTypes/services";
  import { AWS_EXTERNAL_SERVICES, AWS_INTERNAL_SERVICES, AWS_SERVICES } from "$src/helpers/constants";
  import { createEventDispatcher } from "svelte";
  import ApigatewayV2 from "./apigateway/apigatewayV2.svelte";
  import Cdn from "./cloudfront/cdn.svelte";
  import Dynamodb from "./dynamodb/dynamodb.svelte";
  import Ec2 from "./ec2/ec2.svelte";
  import Efs from "./efs/efs.svelte";
  import Eks from "./eks/eks.svelte";
  import ElbV2 from "./elasticloadbalancer/elbV2.svelte";
  import Lambda from "./lambda/lambda.svelte";
  import Rds from "./rds/rds.svelte";
  import Route53 from "./route53/route53.svelte";
  import S3Bucket from "./s3/s3Bucket.svelte";
  import type { HighLightProps, LegendProps, TargetFromNodeProps } from "$src/customTypes/konva";
  import GroupIndex from "../views/groupIndex.svelte";
  import { reorderAwsServices } from "$src/helpers";
  import { getResourceData } from "$src/api/services";
  import { getApiGatewayIntegrations } from "$src/api/aws";
  import { stripIntegrationUriIp } from "$src/helpers/konva/index";
  import { precomputeBorder } from "$src/helpers/aws";

  export let result: ResourceDataReturnType;
  export let projectId: string;
  export let providerAccountId: string;
  export let region: string;
  export let setLoading: (bool: boolean) => void;
  export let setLegend: (data: LegendProps[]) => void;
  export let highlights: HighLightProps;

  let remount = true;

  export let setLineTargets: (tg: TargetFromNodeProps[]) => void;

  const componentMap = {
  	[AWS_SERVICES.EC2]: Ec2,
  	[AWS_SERVICES.APIGATEWAYV2]: ApigatewayV2,
  	[AWS_SERVICES.LAMBDA]: Lambda,
  	[AWS_SERVICES.CLOUDFRONT]: Cdn,
  	[AWS_SERVICES.S3]: S3Bucket,
  	[AWS_SERVICES.RDS]: Rds,
  	[AWS_SERVICES.EFS]: Efs,
  	[AWS_SERVICES.EKS]: Eks,
  	[AWS_SERVICES.ELBV2]: ElbV2,
  	[AWS_SERVICES.ROUTE53]: Route53,
  	[AWS_SERVICES.DYNAMODB]: Dynamodb,
  };

  const dispatch = createEventDispatcher();

  // Group services by external / internal accessibility.
  const externalServices = [
  	AWS_SERVICES.APIGATEWAYV2,
  	AWS_SERVICES.CLOUDFRONT,
  	AWS_SERVICES.ROUTE53,
  	AWS_SERVICES.ELBV2,
  ] as string[];

  const groupResult = (result: ResourceDataReturnType) => {
  	const data = result.reduce(
  		(acc, r) => {
  			if (externalServices.includes(r.name)) {
  				acc["externalGroup"] = (acc["externalGroup"] || []).concat(r);
  				return acc;
  			}
  			acc["internalGroup"] = (acc["internalGroup"] || []).concat(r);
  			return acc;
  		},
      {
      	externalGroup: [],
      	internalGroup: [],
      } as GroupedData
  	);
  	data.externalGroup = reorderAwsServices(data.externalGroup, AWS_EXTERNAL_SERVICES);
  	data.internalGroup = reorderAwsServices(data.internalGroup, AWS_INTERNAL_SERVICES);
  	return data;
  };

  // Precomputing dimensions to draw border


  let groupedData = groupResult(result);

  precomputeBorder(result);

  let groupView: any = null;
  export const updateBorder = () => groupView?.updateBorder();

  export const refetchData = async (region: string) => {
  	try {
  		setLoading(true);
  		remount = false;
  		const resp = await getResourceData({
  			projectId,
  			providerAccountId,
  			region,
  		});
  		if (resp.error || !resp.data) throw resp;
  		const res = resp.data as any;
  		const idex = res.findIndex((r: any) => r.name === AWS_SERVICES.APIGATEWAYV2);
  		if (idex >= 0) {
  			const apigateway = res[idex];
  			// Fetch integrations and attach it to apigateway data
  			// This also makes it easier and faster to build arrow connectors
  			// to show on canvas
  			const apiGatewayWithIntegrations = await Promise.all(
  				apigateway.result.Items.map((item: ApiGatewayV2Props["Items"][0]) => {
  					return getApiGatewayIntegrations({
  						projectId,
  						providerAccountId,
  						region,
  						apiId: item.ApiId,
  					}).then((res) => {
  						return {
  							...item,
  							integrations: res.data?.Items || [],

  							/**
                 * Strip the integrationUri to get the ip address
                 * of Ec2 VMs and lambda functions.
                 *
                 * We will be using this as our target to draw connecting arrows.
                 */
  							lineTargets: (res.data?.Items || []).map((it) => ({
  								from: item.ApiId,
  								to: stripIntegrationUriIp(it.IntegrationUri),
  							})),
  						};
  					});
  				})
  			);
  			res[idex].result =
          apiGatewayWithIntegrations as ApiGatewayWithIntegrationProps[];
  		}
  		groupedData = groupResult(res);
  		precomputeBorder(res);
  		dispatch("remount");
  	} catch (err) {
  		console.error("Unable to fetch data", err);
  	}
  	setLoading(false);
  	remount = true;
  };

  const _result = (item: ResourceDataReturnType[0]) => item.result as any;
</script>

{#if remount}
  <GroupIndex
    {groupedData}
    bind:this={groupView}
  >
    <svelte:fragment let:item let:index let:idx let:externalGroup>
      {#if componentMap[item.name]}
        <svelte:component
          this={componentMap[item.name]}
          {setLegend}
          {idx}
          {externalGroup}
          data={_result(item)}
          {projectId}
          {providerAccountId}
          {region}
		      {highlights}
          on:mouseenter={(e) => {
          	dispatch("mouseenter", e.detail);
          }}
          on:mouseleave={() => {
          	dispatch("mouseleave");
          }}
          on:click={(e) => {
          	dispatch("click", e.detail);
          }}
          {setLineTargets}
          on:initialPosition={(e) => {
          	/**
             * It is recommended to debounce this data before
             * updating it to store to avoid too many update calls.
             */
          	dispatch("initialPosition", e.detail);
          }}
          on:dragmove={(e) => {
          	// has to be updated real-time for smoother experience
          	dispatch("dragmove", e.detail);
          }}
          on:dragend={(e) => {
          	dispatch("dragend", e.detail);
          }}
        />
      {/if}
    </svelte:fragment>
  </GroupIndex>
{/if}

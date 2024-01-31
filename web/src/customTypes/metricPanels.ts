export type MetricPanelProps = {
  id: string;
  name: string;
  description?: string;
  panels: any;
  providerAccountId: string;
  instanceId: string;
  healthStatus: "DISCONNECTED" | "ACTIVE";
  metadata?: Record<string, unknown>;
};

export type CreateMetricPanelProps = Pick<
  MetricPanelProps,
  "name" | "description" | "providerAccountId" | "instanceId" | "panels" | "healthStatus"
>;

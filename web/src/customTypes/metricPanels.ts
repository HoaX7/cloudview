export type PanelProps = {
    type: string;
    name: string;
}

export type MetricPanelProps = {
  id: string;
  name: string;
  description?: string;
  panels: PanelProps[];
  providerAccountId: string;
  instanceId: string;
  healthStatus: "DISCONNECTED" | "ACTIVE";
  metadata?: Record<string, unknown>;
};

export type CreateMetricPanelProps = Pick<
  MetricPanelProps,
  "name" | "description" | "providerAccountId" | "instanceId" | "panels" | "healthStatus"
>;

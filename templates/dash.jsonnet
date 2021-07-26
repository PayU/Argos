local grafana = import '/Users/jacobgo/Argos/vendor/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local template = grafana.template;
local singlestat = grafana.singlestat;
local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

local slo =
  graphPanel.new(
    title='slo',
    datasource='Prometheus',
    linewidth=2,
    format='Bps',
    aliasColors={
      Rx: 'light-green',
      Tx: 'light-red',
    },
  ).addTarget(
    prometheus.target(
      'histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{app="services-iam-api", env="eks-mars-apps"}[1h])) by (le, code, route))',
      legendFormat='Rx',
    )
  );

dashboard.new(
  'Prometheus test',
  tags=['prometheus'],
  schemaVersion=18,
  editable=true,
  time_from='now-1h',
  refresh='1m',
)
.addTemplate(
  template.datasource(
    'PROMETHEUS_DS',
    'prometheus',
    'Prometheus',
    hide='label',
  )
)
.addTemplate(
  template.new(
    'instance',
    '$PROMETHEUS_DS',
    'label_values(node_network_receive_bytes_total, instance)',
    label='Instance',
    refresh='time',
  )
).addTemplate(
  template.new(
    'prom_instance',
    '$PROMETHEUS_DS',
    'label_values(prometheus_build_info, instance)',
    label='Prom Instance',
    refresh='time',   
  )
)


.addPanels(
  [
    slo { gridPos: { h: 4, w: 3, x: 0, y: 0 } }
  ]
)
local prometheus = import "prometheus-ksonnet/prometheus-ksonnet.libsonnet";

prometheus {
  _config+:: {
    cluster_name: "scenarios",
    namespace: "observability",
  },
}

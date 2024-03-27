package config

type Listener struct {
	RoutingUrl        configUrl   `json:"url"`
	TargetGroupConfig TargetGroup `json:"targetGroup"`
}

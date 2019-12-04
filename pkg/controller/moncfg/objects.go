package moncfg

type global struct {
	ResolveTimeout string `yaml:"resolve_timeout"`
}

type routes struct {
	Receiver string            `yaml:"receiver"`
	MatchRe  map[string]string `yaml:"match_re"`
}

type route struct {
	GroupBy        []string `yaml:"group_by"`
	GroupWait      string   `yaml:"group_wait"`
	GroupInterval  string   `yaml:"group_interval"`
	RepeatInterval string   `yaml:"repeat_interval"`
	Receiver       string   `yaml:"receiver"`
	Routes         []routes `yaml:"routes,omitempty"`
}

type slackconfig struct {
	ApiURL  string `yaml:"api_url"`
	Channel string `yaml:"channel"`
}

type emailconfig struct {
	To           string `yaml:"to"`
	From         string `yaml:"from"`
	SmartHost    string `yaml:"smarthost"`
	AuthUsername string `yaml:"auth_username"`
	AuthIdentity string `yaml:"auth_identity"`
	AuthPassword string `yaml:"auth_password"`
}

type receiver struct {
	Name         string        `yaml:"name"`
	SlackConfigs []slackconfig `yaml:"slack_configs,omitempty"`
	EmailConfigs []emailconfig `yaml:"email_configs,omitempty"`
}

type alertConfig struct {
	Global    global     `yaml:"global"`
	Route     route      `yaml:"route"`
	Receivers []receiver `yaml:"receivers"`
}

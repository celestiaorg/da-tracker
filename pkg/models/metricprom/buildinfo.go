package metricprom

type BuildInfo struct {
	Metric struct {
		Name            string `json:"__name__"`
		BuildTime       string `json:"build_time"`
		GoLangVersion   string `json:"golang_version"`
		Instance        string `json:"instance"`
		Job             string `json:"job"`
		LastCommit      string `json:"last_commit"`
		SemanticVersion string `json:"semantic_version"`
		SystemVersion   string `json:"system_version"`
	} `json:"metric"`
	Value []interface{} `json:"value"`
}

type BuildInfoResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string      `json:"resultType"`
		Result     []BuildInfo `json:"result"`
	} `json:"data"`
}

type BuildInfoDetails struct {
	SemanticVersion string
	Job             string
	SystemVersion   string
}

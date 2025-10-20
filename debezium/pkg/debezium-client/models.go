package debezium_client

type GetConnectorsResponse struct {
	Connectors []struct {
		Status struct {
			Name      string `json:"name"`
			Connector struct {
				State    string `json:"state"`
				WorkerId string `json:"worker_id"`
			} `json:"connector"`
			Tasks []struct {
				Id       int    `json:"id"`
				State    string `json:"state"`
				WorkerId string `json:"worker_id"`
			} `json:"tasks"`
			Type string `json:"type"`
		} `json:"status"`
	} `json:"connector"`
}

type CreateConnectorRequest struct {
	Name   string                `json:"name"`
	Config CreateConnectorConfig `json:"config"`
}

type CreateConnectorConfig struct {
	ConnectorClass       string            `json:"connector.class"`
	TasksMax             string            `json:"tasks.max"`
	DatabaseHostname     string            `json:"database.hostname"`
	DatabasePort         string            `json:"database.port"`
	DatabaseUser         string            `json:"database.user"`
	DatabasePassword     string            `json:"database.password"`
	DatabaseDbname       string            `json:"database.dbname"`
	DatabaseServerName   string            `json:"database.server.name"`
	AdditionalParameters map[string]string `json:"-,omitempty"`
}

type CreateConnectorResponse struct {
	Name   string                `json:"name"`
	Config CreateConnectorConfig `json:"config"`
	Tasks  []any                 `json:"tasks"`
	Type   string                `json:"type"`
}

type GetConnectorResponse struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	Tasks  []TaskInfo             `json:"tasks"`
	Type   string                 `json:"type"`
}

type TaskInfo struct {
	Connector string         `json:"connector"`
	Task      int            `json:"task"`
	Config    map[string]any `json:"config"`
}

type ConnectorStatus struct {
	Name      string         `json:"name"`
	Connector ConnectorState `json:"connector"`
	Tasks     []TaskState    `json:"tasks"`
	Type      string         `json:"type"`
}

type ConnectorState struct {
	State    string `json:"state"`
	WorkerId string `json:"worker_id"`
}

type TaskState struct {
	Id       int    `json:"id"`
	State    string `json:"state"`
	WorkerId string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
}

type UpdateConnectorConfigRequest struct {
	ConnectorClass       string            `json:"connector.class"`
	TasksMax             string            `json:"tasks.max"`
	AdditionalParameters map[string]string `json:"-,omitempty"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type ListConnectorsResponse struct {
	Names    []string
	Statuses map[string]ConnectorStatus
}

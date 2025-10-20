package debezium_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	postCreateConnectors  = "/connectors"
	getConnector          = "/connectors/%s"
	deleteConnector       = "/connectors/%s"
	getConnectorStatus    = "/connectors/%s/status"
	updateConnectorConfig = "/connectors/%s/config"
	pauseConnector        = "/connectors/%s/pause"
	resumeConnector       = "/connectors/%s/resume"
	restartConnector      = "/connectors/%s/restart"
	getConnectorTasks     = "/connectors/%s/tasks"
	restartConnectorTask  = "/connectors/%s/tasks/%d/restart"
	listConnectors        = "/connectors"
)

var (
	ErrEmptyConnectorName = errors.New("connector name cannot be empty")
)

func validateConnectorName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyConnectorName
	}
	return nil
}

func (c *Client) CreateConnector(ctx context.Context, data CreateConnectorRequest) (*CreateConnectorResponse, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("CreateConnector.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+postCreateConnectors, bytes.NewBuffer(d))
	if err != nil {
		return nil, fmt.Errorf("CreateConnector.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CreateConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("CreateConnector.DecodeError: %w", err)
		}
		return nil, fmt.Errorf("CreateConnector: %s", errorResponse.Message)
	}

	var connectorResponse CreateConnectorResponse
	if err := json.NewDecoder(resp.Body).Decode(&connectorResponse); err != nil {
		return nil, fmt.Errorf("CreateConnector.UnmarshalJSON: %w", err)
	}

	return &connectorResponse, nil
}

func (c *Client) GetConnector(ctx context.Context, name string) (GetConnectorResponse, error) {
	if err := validateConnectorName(name); err != nil {
		return GetConnectorResponse{}, err
	}

	var connectorResponse GetConnectorResponse

	url := fmt.Sprintf(c.baseURL+getConnector, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return GetConnectorResponse{}, fmt.Errorf("GetConnector.DecodeError: %w", err)
		}
		return GetConnectorResponse{}, fmt.Errorf("GetConnector: %s", errorResponse.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&connectorResponse); err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.UnmarshalJSON: %w", err)
	}

	return connectorResponse, nil
}

func (c *Client) GetConnectorStatus(ctx context.Context, name string) (ConnectorStatus, error) {
	if err := validateConnectorName(name); err != nil {
		return ConnectorStatus{}, err
	}

	var status ConnectorStatus

	url := fmt.Sprintf(c.baseURL+getConnectorStatus, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ConnectorStatus{}, fmt.Errorf("GetConnectorStatus.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return ConnectorStatus{}, fmt.Errorf("GetConnectorStatus.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return ConnectorStatus{}, fmt.Errorf("GetConnectorStatus.DecodeError: %w", err)
		}
		return ConnectorStatus{}, fmt.Errorf("GetConnectorStatus: %s", errorResponse.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return ConnectorStatus{}, fmt.Errorf("GetConnectorStatus.UnmarshalJSON: %w", err)
	}

	return status, nil
}

func (c *Client) DeleteConnector(ctx context.Context, name string) error {
	if err := validateConnectorName(name); err != nil {
		return err
	}

	url := fmt.Sprintf(c.baseURL+deleteConnector, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("DeleteConnector.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return fmt.Errorf("DeleteConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("DeleteConnector.DecodeError: %w", err)
		}
		return fmt.Errorf("DeleteConnector: %s", errorResponse.Message)
	}

	return nil
}

func (c *Client) UpdateConnectorConfig(ctx context.Context, name string, config map[string]interface{}) (GetConnectorResponse, error) {
	if err := validateConnectorName(name); err != nil {
		return GetConnectorResponse{}, err
	}

	var connectorResponse GetConnectorResponse

	data, err := json.Marshal(config)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig.Marshal: %w", err)
	}

	url := fmt.Sprintf(c.baseURL+updateConnectorConfig, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cc.Do(req)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig.DecodeError: %w", err)
		}
		return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig: %s", errorResponse.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&connectorResponse); err != nil {
		return GetConnectorResponse{}, fmt.Errorf("UpdateConnectorConfig.UnmarshalJSON: %w", err)
	}

	return connectorResponse, nil
}

func (c *Client) PauseConnector(ctx context.Context, name string) error {
	if err := validateConnectorName(name); err != nil {
		return err
	}

	url := fmt.Sprintf(c.baseURL+pauseConnector, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("PauseConnector.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return fmt.Errorf("PauseConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("PauseConnector.DecodeError: %w", err)
		}
		return fmt.Errorf("PauseConnector: %s", errorResponse.Message)
	}

	return nil
}

func (c *Client) ResumeConnector(ctx context.Context, name string) error {
	if err := validateConnectorName(name); err != nil {
		return err
	}

	url := fmt.Sprintf(c.baseURL+resumeConnector, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("ResumeConnector.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return fmt.Errorf("ResumeConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("ResumeConnector.DecodeError: %w", err)
		}
		return fmt.Errorf("ResumeConnector: %s", errorResponse.Message)
	}

	return nil
}

func (c *Client) RestartConnector(ctx context.Context, name string) error {
	if err := validateConnectorName(name); err != nil {
		return err
	}

	url := fmt.Sprintf(c.baseURL+restartConnector, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("RestartConnector.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return fmt.Errorf("RestartConnector.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("RestartConnector.DecodeError: %w", err)
		}
		return fmt.Errorf("RestartConnector: %s", errorResponse.Message)
	}

	return nil
}

func (c *Client) GetConnectorTasks(ctx context.Context, name string) ([]TaskInfo, error) {
	if err := validateConnectorName(name); err != nil {
		return nil, err
	}

	var tasks []TaskInfo

	url := fmt.Sprintf(c.baseURL+getConnectorTasks, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetConnectorTasks.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetConnectorTasks.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("GetConnectorTasks.DecodeError: %w", err)
		}
		return nil, fmt.Errorf("GetConnectorTasks: %s", errorResponse.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("GetConnectorTasks.UnmarshalJSON: %w", err)
	}

	return tasks, nil
}

func (c *Client) RestartConnectorTask(ctx context.Context, name string, taskId int) error {
	if err := validateConnectorName(name); err != nil {
		return err
	}

	url := fmt.Sprintf(c.baseURL+restartConnectorTask, name, taskId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("RestartConnectorTask.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return fmt.Errorf("RestartConnectorTask.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("RestartConnectorTask.DecodeError: %w", err)
		}
		return fmt.Errorf("RestartConnectorTask: %s", errorResponse.Message)
	}

	return nil
}

func (c *Client) ListConnectors(ctx context.Context, expandStatus bool) (ListConnectorsResponse, error) {
	var result ListConnectorsResponse

	endpoint := c.baseURL + listConnectors
	if expandStatus {
		endpoint += "?expand=status"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return ListConnectorsResponse{}, fmt.Errorf("ListConnectors.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return ListConnectorsResponse{}, fmt.Errorf("ListConnectors.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return ListConnectorsResponse{}, fmt.Errorf("ListConnectors.DecodeError: %w", err)
		}
		return ListConnectorsResponse{}, fmt.Errorf("ListConnectors: %s", errorResponse.Message)
	}

	if expandStatus {
		var statusMap map[string]ConnectorStatus
		if err := json.NewDecoder(resp.Body).Decode(&statusMap); err != nil {
			return ListConnectorsResponse{}, fmt.Errorf("ListConnectors.UnmarshalJSON: %w", err)
		}

		result.Statuses = statusMap
		result.Names = make([]string, 0, len(statusMap))
		for name := range statusMap {
			result.Names = append(result.Names, name)
		}
	} else {
		var names []string
		if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
			return ListConnectorsResponse{}, fmt.Errorf("ListConnectors.UnmarshalJSON: %w", err)
		}
		result.Names = names
	}

	return result, nil
}

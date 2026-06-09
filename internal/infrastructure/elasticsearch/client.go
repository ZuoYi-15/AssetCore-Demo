package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"asset-core/internal/config"
	"asset-core/internal/module/asset"
)

type Client struct {
	enabled    bool
	address    string
	username   string
	password   string
	index      string
	httpClient *http.Client
}

func New(cfg config.ElasticsearchConfig) *Client {
	return &Client{
		enabled:  cfg.Enabled,
		address:  strings.TrimRight(cfg.Address, "/"),
		username: cfg.Username,
		password: cfg.Password,
		index:    cfg.Index,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled && c.address != "" && c.index != ""
}

func (c *Client) IndexAsset(ctx context.Context, item asset.Asset) error {
	if !c.Enabled() {
		return nil
	}
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/%s/_doc/%d?refresh=false", c.index, item.ID)
	return c.do(ctx, http.MethodPut, path, body, nil)
}

func (c *Client) DeleteAsset(ctx context.Context, id uint64) error {
	if !c.Enabled() {
		return nil
	}
	path := fmt.Sprintf("/%s/_doc/%d?refresh=false", c.index, id)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) SearchAssets(ctx context.Context, q asset.Query, offset, limit int) (*asset.SearchResult, error) {
	if !c.Enabled() {
		return nil, nil
	}
	query := buildAssetQuery(q, offset, limit)
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source asset.Asset `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := c.do(ctx, http.MethodPost, "/"+c.index+"/_search", body, &raw); err != nil {
		return nil, err
	}
	items := make([]asset.Asset, 0, len(raw.Hits.Hits))
	for _, hit := range raw.Hits.Hits {
		items = append(items, hit.Source)
	}
	return &asset.SearchResult{Items: items, Total: raw.Hits.Total.Value}, nil
}

func buildAssetQuery(q asset.Query, offset, limit int) map[string]interface{} {
	must := make([]interface{}, 0)
	filter := make([]interface{}, 0)
	if q.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q.Keyword,
				"fields": []string{"asset_name^3", "serial_number^2", "identity_id^2", "mac_address", "ip_address", "hostname", "vendor", "model", "owner_user", "owner_department", "location"},
			},
		})
	}
	if q.Status != "" {
		filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"status.keyword": q.Status}})
	}
	if q.AssetType != "" {
		filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"asset_type.keyword": q.AssetType}})
	}
	if len(must) == 0 {
		must = append(must, map[string]interface{}{"match_all": map[string]interface{}{}})
	}
	return map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"sort": []interface{}{map[string]interface{}{"id": map[string]interface{}{"order": "desc"}}},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body []byte, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, c.address+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("elasticsearch %s %s failed: status=%d body=%s", method, path, resp.StatusCode, string(payload))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

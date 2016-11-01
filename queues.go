package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Information about backing queue (queue storage engine).
type BackingQueueStatus struct {
	Q1                    int     `json:"q1"`
	Q2                    int     `json:"q2"`
	Q3                    int     `json:"q3"`
	Q4                    int     `json:"q4"`
	Length                int64   `json:"len"`
	PendingAcks           int64   `json:"pending_acks"`
	RAMMessageCount       int64   `json:"ram_msg_count"`
	RAMAckCount           int64   `json:"ram_ack_count"`
	PersistentCount       int64   `json:"persistent_count"`
	AverageIngressRate    float64 `json:"avg_ingress_rate"`
	AverageEgressRate     float64 `json:"avg_egress_rate"`
	AverageAckIngressRate float32 `json:"avg_ack_ingress_rate"`
	AverageAckEgressRate  float32 `json:"avg_ack_egress_rate"`
}

type OwnerPidDetails struct {
	Name     string `json:"name"`
	PeerPort Port   `json:"peer_port"`
	PeerHost string `json:"peer_host"`
}

type QueueInfo struct {
	Name                          string                 `json:"name"`
	Vhost                         string                 `json:"vhost"`
	Durable                       bool                   `json:"durable"`
	AutoDelete                    bool                   `json:"auto_delete"`
	Arguments                     map[string]interface{} `json:"arguments"`
	Node                          string                 `json:"node"`
	Status                        string                 `json:"status"`
	Memory                        int64                  `json:"memory"`
	Consumers                     int                    `json:"consumers"`
	ConsumerUtilisation           interface{}            `json:"consumer_utilisation"`
	ExclusiveConsumerTag          interface{}            `json:"exclusive_consumer_tag"`
	Policy                        string                 `json:"policy"`
	MessagesBytes                 int64                  `json:"message_bytes"`
	MessagesBytesPersistent       int64                  `json:"message_bytes_persistent"`
	MessagesBytesRAM              int64                  `json:"message_bytes_ram"`
	Messages                      int                    `json:"messages"`
	MessagesDetails               RateDetails            `json:"messages_details"`
	MessagesPersistent            int                    `json:"messages_persistent"`
	MessagesRAM                   int                    `json:"messages_ram"`
	MessagesReady                 int                    `json:"messages_ready"`
	MessagesReadyDetails          RateDetails            `json:"messages_ready_details"`
	MessagesUnacknowledged        int                    `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails RateDetails            `json:"messages_unacknowledged_details"`
	MessageStats                  MessageStats           `json:"message_stats"`
	OwnerPidDetails               OwnerPidDetails        `json:"owner_pid_details"`
	BackingQueueStatus            BackingQueueStatus     `json:"backing_queue_status"`
	ConsumerDetails               []*ConsumerDetails     `json:"consumer_details,omitempty"`
	State                         string                 `json:"state"`
}

type ConsumerDetails struct {
	ChannleDetails *ChannelDetails `json:"channel_details,omitempty"`
	Queue          *Queue          `json:"queue,omitempty"`
}

type ChannelDetails struct {
	Name           string `json:"name"`
	Number         int    `json:"number,int"`
	User           string `json:"user"`
	ConnectionName string `json:"connection_name"`
	PeerPort       int    `json:"peer_port,int"`
	PeerHost       string `json:"peer_host"`
}

type Queue struct {
	Name  string `json:"name"`
	Vhost string `json:"vhost"`
}

type DetailedQueueInfo QueueInfo

//
// GET /api/queues
//

func (c *Client) ListQueues() (rec []QueueInfo, err error) {
	req, err := newGETRequest(c, "queues")
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

//
// GET /api/queues/{vhost}
//

func (c *Client) ListQueuesIn(vhost string) (rec []QueueInfo, err error) {
	req, err := newGETRequest(c, "queues/"+url.QueryEscape(vhost))
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

//
// GET /api/queues/{vhost}/{name}
//

func (c *Client) GetQueue(vhost, queue string) (rec *DetailedQueueInfo, err error) {
	req, err := newGETRequest(c, "queues/"+url.QueryEscape(vhost)+"/"+queue)
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/exchanges/{vhost}/{exchange}
//

type QueueSettings struct {
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete"`
	Arguments  map[string]interface{} `json:"arguments"`
}

func (c *Client) DeclareQueue(vhost, queue string, info QueueSettings) (res *http.Response, err error) {
	if info.Arguments == nil {
		info.Arguments = make(map[string]interface{})
	}
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "queues/"+url.QueryEscape(vhost)+"/"+url.QueryEscape(queue), body)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/queues/{vhost}/{name}
//

func (c *Client) DeleteQueue(vhost, queue string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "queues/"+url.QueryEscape(vhost)+"/"+url.QueryEscape(queue), nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/queues/{vhost}/{name}/contents
//

func (c *Client) PurgeQueue(vhost, queue string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "queues/"+url.QueryEscape(vhost)+"/"+url.QueryEscape(queue)+"/contents", nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

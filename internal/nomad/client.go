package nomad

import (
	"io"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomadda/internal/structs"
)

type Client struct {
	nomadClient *api.Client
}

type Config struct {
	Address string
	Region  string
	Token   string

	ScanInterval time.Duration
}

func NewClient(c *Config) (*Client, error) {
	cfg := api.DefaultConfig()
	if c.Address != "" {
		cfg.Address = c.Address
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		nomadClient: client,
	}, nil
}

func (c *Client) NomadInfo() (structs.NomadInfo, error) {
	members, err := c.nomadClient.Agent().Members()
	if err != nil {
		return structs.NomadInfo{}, err
	}
	servers := []string{}
	for _, m := range members.Members {
		servers = append(servers, m.Name)
	}
	sort.Strings(servers)

	nodeList, _, err := c.nomadClient.Nodes().List(&api.QueryOptions{})
	if err != nil {
		return structs.NomadInfo{}, err
	}

	nodes := []string{}
	for _, n := range nodeList {
		nodes = append(nodes, n.Name)
	}
	sort.Strings(nodes)

	return structs.NomadInfo{
		Servers: structs.Servers(servers),
		Clients: structs.Clients(nodes),
	}, nil
}

func (c *Client) Logs() (io.ReadCloser, error) {
	cancel := make(chan struct{})

	allocList, _, err := c.nomadClient.Allocations().List(&api.QueryOptions{})
	if err != nil {
		spew.Dump(len(allocList))
		return nil, err
	}

	allocs := []*api.Allocation{}
	for _, a := range allocList {
		alloc, _, err := c.nomadClient.Allocations().Info(a.ID, &api.QueryOptions{})
		if err != nil {
			return nil, err
		}

		allocs = append(allocs, alloc)
	}

	tasks := allocs[0].GetTaskGroup().Tasks

	frames, errCh := c.nomadClient.AllocFS().Logs(allocs[0], true, tasks[0].Name, "stdout", api.OriginEnd, 0, cancel, &api.QueryOptions{})
	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	var r io.ReadCloser
	frameReader := api.NewFrameReader(frames, errCh, cancel)
	frameReader.SetUnblockTime(500 * time.Millisecond)
	r = frameReader
	go func() {
		<-signalCh
		// eventually handle this gracefully
		r.Close()
	}()
	return r, nil
}

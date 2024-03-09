package c2

import (
	"fmt"
	"sync"
	"time"
)

type Agent struct {
	Id string
	Ip string
	LastCall time.Time
	CmdQueue [][]string
}

type SafeAgentMap struct {
	mtx sync.Mutex
	Agents map[string]*Agent // Key: Agent id, Value: Pointer to Agent
}

var AgentMap SafeAgentMap = SafeAgentMap{Agents: make(map[string]*Agent)}

func (am *SafeAgentMap) Add(agent *Agent) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	// Check if agent is in the map
	if _, exists := am.Agents[agent.Id]; !exists {
		am.Agents[agent.Id] = agent
	}
}

func (am *SafeAgentMap) Get(agentId string) *Agent {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if agent, exists := am.Agents[agentId]; exists {
		return agent
	}
	return nil
}

func (am *SafeAgentMap) Enqueue(agentId string, cmd []string) error {
	agent := am.Get(agentId)
	if agent == nil {
		return fmt.Errorf("agent '%s' does not exist.", agentId)
	}
	am.mtx.Lock()
	defer am.mtx.Unlock()
	agent.CmdQueue = append(agent.CmdQueue, cmd)
	return nil
}

func (am *SafeAgentMap) Dequeue(agentId string) ([]string, error) {
	agent := am.Get(agentId)
	if agent == nil {
		return nil, fmt.Errorf("agent '%s' does not exist.", agentId)
	}
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if len(agent.CmdQueue) < 1 {
		return nil, fmt.Errorf("agent '%s' doesnt have any queued commands", agentId)
	}
	cmd := agent.CmdQueue[0]
	agent.CmdQueue = agent.CmdQueue[1:]
	return cmd, nil // nil cos no errors
}
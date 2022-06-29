// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// util
package infoutil

import (
	"fmt"
	"net"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/marmotedu/errors"
	"github.com/notone0010/pigpig/pkg/log"
)

func GetTransportAddr() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip, nil
}

var distributedId *DistributedId

type DistributedId struct {
	Node      *snowflake.Node
	MachineId int64
}

func NewDistributedId(machineId int64) (*DistributedId, error) {
	node, err := snowflake.NewNode(machineId)
	if err != nil {
		return nil, err
	}
	distributedId = &DistributedId{
		MachineId: machineId,
		Node:      node,
	}
	return distributedId, nil
}

func (d *DistributedId) Generate() snowflake.ID {
	return d.Node.Generate()
}

func (d *DistributedId) ParseId(id interface{}) (snowflake.ID, error) {
	switch id {
	case id.(int64):
		return snowflake.ParseInt64(id.(int64)), nil
	case id.(string):
		return snowflake.ParseString(id.(string))
	default:
		return 0, errors.New("id cannot transport known types")
	}
}

func GetDistributedNode() *DistributedId {
	if distributedId == nil {
		distributedId, err := NewDistributedId(1)
		if err != nil {
			log.Fatalf("failed to create DistributedId --- %s", err.Error())
			return nil
		}
		return distributedId
	}
	return distributedId
}

func SetDistributedNode(d *DistributedId) {
	distributedId = d
}

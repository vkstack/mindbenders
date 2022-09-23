package packet

import (
	"context"
	"encoding/json"

	"gitlab.com/dotpe/mindbenders/corel"
)

type Packet interface {
	Ctx() context.Context
	CorelId() *corel.CoRelationId
}

type packet struct {
	corelid *corel.CoRelationId
	ctx     context.Context

	Correlation string      `json:"correlation"`
	Payload     interface{} `json:"data"`
}

func (p *packet) Ctx() context.Context {
	return p.ctx
}

func (p *packet) CorelId() *corel.CoRelationId {
	return p.corelid
}

func NewPacket(raw []byte, dst interface{}) (Packet, error) {
	p := packet{
		Payload: dst,
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	p.corelid = corel.DecodeCorelationId(p.Correlation).Sibling()
	p.ctx = corel.NewCorelCtxFromCorel(p.corelid)
	return &p, nil
}

func NewPacketFromStr(raw string, dst interface{}) (Packet, error) {
	return NewPacket([]byte(raw), dst)
}

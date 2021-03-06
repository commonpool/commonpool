package commands

import (
	"fmt"
)

type CommandMapperFunc func(commandType string, bytes []byte) (Command, error)

type CommandMapper struct {
	mapperFuncs map[string]CommandMapperFunc
}

func NewCommandMapper() *CommandMapper {
	return &CommandMapper{
		mapperFuncs: map[string]CommandMapperFunc{},
	}
}

func (c *CommandMapper) RegisterMapper(commandType string, mapper CommandMapperFunc) {
	c.mapperFuncs[commandType] = mapper
}

func (c *CommandMapper) Map(commandType string, command []byte) (Command, error) {
	mapper, ok := c.mapperFuncs[commandType]
	if !ok {
		return nil, fmt.Errorf("no mapper for command %s", commandType)
	}
	return mapper(commandType, command)
}

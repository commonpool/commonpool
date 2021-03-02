package service

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type tradingTestSuite struct {
	suite.Suite
}

func TestTrading(t *testing.T) {
	suite.Run(t, &tradingTestSuite{})
}

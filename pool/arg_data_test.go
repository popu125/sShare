package pool

import (
	"encoding/json"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

const port = 2596

type MSuite struct {
	argData *runArgData
	pass    string
	uuid    string
}

var _ = Suite(&MSuite{})

func (s *MSuite) SetUpSuite(c *C) {
	s.argData = NewArgData(port)
	s.pass = s.argData.Pass("pass")
	s.uuid = s.argData.UUID("uuid")
}

func (s *MSuite) TestNewArgData(c *C) {
	c.Assert(s.argData.Port, Equals, port)
}

func (s *MSuite) TestRunArgData_Pass(c *C) {
	c.Assert(s.pass, Equals, s.argData.Pass("pass"))
	c.Assert(s.pass, Not(Equals), s.argData.Pass("pass2"))
}

func (s *MSuite) TestRunArgData_UUID(c *C) {
	c.Assert(s.uuid, Equals, s.argData.UUID("uuid"))
	c.Assert(s.uuid, Not(Equals), s.argData.UUID("uuid2"))
}

func (s *MSuite) TestRunArgData_Data(c *C) {
	dataStr := s.argData.Data()
	c.Assert(dataStr, Not(Equals), "")
	obj := &runArgData{}
	json.Unmarshal([]byte(dataStr), obj)
	c.Assert(obj.Port, Equals, port)
	c.Assert(obj.PassMap["pass"], Equals, s.pass)
	c.Assert(obj.UUIDMap["uuid"], Equals, s.uuid)
}

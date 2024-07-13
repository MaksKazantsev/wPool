package tests

import (
	"errors"
	"github.com/MaksKazantsev/wPool"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(PoolTestSuite))
}

type PoolTestSuite struct {
	suite.Suite
	pool wPool.Pool
}

func (p *PoolTestSuite) SetupTest() {
	p.pool = wPool.NewPool(3)
}

func (p *PoolTestSuite) TestBasic() {
	err := errors.New("worker error")

	tasks := []struct {
		fn  wPool.Tasker
		err error
	}{
		{
			fn: func() error {
				return err
			},
			err: err,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return err
			},
			err: err,
		},
		{
			fn: func() error {
				return err
			},
			err: err,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},

		{
			fn: func() error {
				p.pool.Stop()
				return nil
			},
			err: nil,
		},
	}

	errAmount := 0

	for _, t := range tasks {
		p.pool.Task(t.fn)
		if t.err != nil {
			errAmount++
		}
	}

	for err = range p.pool.CatchError() {
		p.Require().Error(err)
		errAmount--
	}

	p.Require().Equal(errAmount, 0)
}

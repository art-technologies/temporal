// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package vclock

import (
	"fmt"

	clockpb "go.temporal.io/server/api/clock/v1"
	"go.temporal.io/server/common/persistence"
)

type (
	ShardOwnership interface {
		GetVClock() *clockpb.ShardClock
		AssertOwnership() error
	}
)

func AssertConsistency(
	clock *clockpb.ShardClock,
	shardOwnership ShardOwnership,
) error {
	result, err := Compare(clock, shardOwnership.GetVClock())
	if err != nil {
		return err
	}

	// clock happens before | same as shard's clock
	if result <= 0 {
		return nil
	}
	return &persistence.ShardOwnershipLostError{
		ShardID: clock.GetId(),
		Msg:     fmt.Sprintf("Shard %v consistency check failed, reloading", clock.GetId()),
	}
}

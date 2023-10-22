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

package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ Batch = (*batch)(nil)

type (
	batch struct {
		session *session

		mongoBatch *mongo.Batch
	}
)

// Definition of all BatchTypes
const (
	LoggedBatch BatchType = iota
	UnloggedBatch
	CounterBatch
)

func newBatch(
	session *session,
	mongoBatch *mongo.Batch,
) *batch {
	return &batch{
		session:    session,
		mongoBatch: mongoBatch,
	}
}

func (b *batch) Query(stmt string, args ...interface{}) {
	b.mongoBatch.Query(stmt, args...)
}

func (b *batch) WithContext(ctx context.Context) Batch {
	return newBatch(b.session, b.mongoBatch.WithContext(ctx))
}

func (b *batch) WithTimestamp(timestamp int64) Batch {
	b.mongoBatch.WithTimestamp(timestamp)
	return newBatch(b.session, b.mongoBatch)
}

func mustConvertBatchType(batchType BatchType) mongo.BatchType {
	switch batchType {
	case LoggedBatch:
		return mongo.LoggedBatch
	case UnloggedBatch:
		return mongo.UnloggedBatch
	case CounterBatch:
		return mongo.CounterBatch
	default:
		panic(fmt.Sprintf("Unknown mongo BatchType: %v", batchType))
	}
}

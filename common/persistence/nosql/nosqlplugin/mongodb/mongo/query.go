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

	"go.mongodb.org/mongo-driver/mongo"
)

var _ Query = (*query)(nil)

type (
	query struct {
		session    *session
		mongoQuery *mongo.Query
	}
)

func newQuery(
	session *session,
	mongoQuery *mongo.Query,
) *query {
	return &query{
		session:    session,
		mongoQuery: mongoQuery,
	}
}

func (q *query) Exec() (retError error) {
	defer func() { q.session.handleError(retError) }()

	return q.mongoQuery.Exec()
}

func (q *query) Scan(
	dest ...interface{},
) (retError error) {
	defer func() { q.session.handleError(retError) }()

	return q.mongoQuery.Scan(dest...)
}

func (q *query) ScanCAS(
	dest ...interface{},
) (_ bool, retError error) {
	defer func() { q.session.handleError(retError) }()

	return q.mongoQuery.ScanCAS(dest...)
}

func (q *query) MapScan(
	m map[string]interface{},
) (retError error) {
	defer func() { q.session.handleError(retError) }()

	return q.mongoQuery.MapScan(m)
}

func (q *query) MapScanCAS(
	dest map[string]interface{},
) (_ bool, retError error) {
	defer func() { q.session.handleError(retError) }()

	return q.mongoQuery.MapScanCAS(dest)
}

func (q *query) Iter() Iter {
	iter := q.mongoQuery.Iter()
	return newIter(q.session, iter)
}

func (q *query) PageSize(n int) Query {
	q.mongoQuery.PageSize(n)
	return newQuery(q.session, q.mongoQuery)
}

func (q *query) PageState(state []byte) Query {
	q.mongoQuery.PageState(state)
	return newQuery(q.session, q.mongoQuery)
}

func (q *query) Consistency(c Consistency) Query {
	q.mongoQuery.Consistency(mustConvertConsistency(c))
	return newQuery(q.session, q.mongoQuery)
}

func (q *query) WithTimestamp(timestamp int64) Query {
	q.mongoQuery.WithTimestamp(timestamp)
	return newQuery(q.session, q.mongoQuery)
}

func (q *query) WithContext(ctx context.Context) Query {
	q2 := q.mongoQuery.WithContext(ctx)
	if q2 == nil {
		return nil
	}
	return newQuery(q.session, q2)
}

func (q *query) Bind(v ...interface{}) Query {
	q.mongoQuery.Bind(v...)
	return newQuery(q.session, q.mongoQuery)
}

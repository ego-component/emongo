// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package emongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type processor func(fn processFn) error
type processFn func(*cmd) error

type cmd struct {
	name     string
	req      []interface{}
	res      interface{}
	dbName   string
	collName string
}

func logCmd(logMode bool, c *cmd, name string, res interface{}, req ...interface{}) {
	// 只有开启log模式才会记录req、res
	if logMode {
		c.name = name
		c.req = append(c.req, req...)
		switch res := res.(type) {
		case *mongo.SingleResult:
			val, _ := res.DecodeBytes()
			c.res = val
		default:
			c.res = res
		}
	}
}

type Collection struct {
	coll      *mongo.Collection
	processor processor
	logMode   bool
}

func (wc *Collection) cmd(c *cmd) *cmd {
	c.dbName = wc.coll.Database().Name()
	c.collName = wc.coll.Name()
	return c
}

func (wc *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (res *mongo.Cursor, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.Aggregate(ctx, pipeline, opts...)
		logCmd(wc.logMode, wc.cmd(c), "Aggregate", res, pipeline)
		return err
	})
	return
}

func (wc *Collection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (
	res *mongo.BulkWriteResult, err error) {

	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.BulkWrite(ctx, models, opts...)
		logCmd(wc.logMode, wc.cmd(c), "BulkWrite", res, models)
		return err
	})
	return
}

func (wc *Collection) Clone(opts ...*options.CollectionOptions) (res *mongo.Collection, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.Clone(opts...)
		logCmd(wc.logMode, wc.cmd(c), "Clone", res)
		return err
	})
	return
}

func (wc *Collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (res int64, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.CountDocuments(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "CountDocuments", res, filter)
		return err
	})
	return res, err
}

func (wc *Collection) Database() *mongo.Database { return wc.coll.Database() }

func (wc *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (
	res *mongo.DeleteResult, err error) {

	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.DeleteMany(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "DeleteMany", res, filter)
		return err
	})
	return
}

func (wc *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (res *mongo.DeleteResult, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.DeleteOne(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "DeleteOne", res, filter)
		return err
	})
	return
}

func (wc *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (res []interface{}, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.Distinct(ctx, fieldName, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "Distinct", nil, fieldName, filter)
		return err
	})
	return
}

func (wc *Collection) Drop(ctx context.Context) error {
	return wc.processor(func(c *cmd) error {
		logCmd(wc.logMode, wc.cmd(c), "Drop", nil)
		return wc.coll.Drop(ctx)
	})
}

func (wc *Collection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (res int64, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.EstimatedDocumentCount(ctx, opts...)
		logCmd(wc.logMode, wc.cmd(c), "EstimatedDocumentCount", res)
		return err
	})
	return
}

func (wc *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (res *mongo.Cursor, err error) {
	err = wc.processor(func(c *cmd) error {
		res, err = wc.coll.Find(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "Find", res, filter)
		return err
	})
	return
}

func (wc *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (res *mongo.SingleResult) {
	_ = wc.processor(func(c *cmd) error {
		res = wc.coll.FindOne(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "FindOne", res, filter)
		return res.Err()
	})
	return
}

func (wc *Collection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (res *mongo.SingleResult) {
	_ = wc.processor(func(c *cmd) error {
		res = wc.coll.FindOneAndDelete(ctx, filter, opts...)
		logCmd(wc.logMode, wc.cmd(c), "FindOneAndDelete", res, filter)
		return res.Err()
	})
	return
}

func (wc *Collection) FindOneAndReplace(ctx context.Context, filter, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) (res *mongo.SingleResult) {
	_ = wc.processor(func(c *cmd) error {
		res = wc.coll.FindOneAndReplace(ctx, filter, replacement, opts...)
		logCmd(wc.logMode, wc.cmd(c), "FindOneAndReplace", res, filter)
		return res.Err()
	})
	return
}

func (wc *Collection) FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (res *mongo.SingleResult) {
	_ = wc.processor(func(c *cmd) error {
		res = wc.coll.FindOneAndUpdate(ctx, filter, update, opts...)
		logCmd(wc.logMode, wc.cmd(c), "FindOneAndReplace", res, filter)
		return res.Err()
	})
	return
}

func (wc *Collection) Indexes() mongo.IndexView { return wc.coll.Indexes() }

func (wc *Collection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (res *mongo.InsertManyResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.InsertMany(ctx, documents, opts...)
		logCmd(wc.logMode, wc.cmd(c), "FindOneAndReplace", res, documents)
		return err
	})
	return
}

func (wc *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (res *mongo.InsertOneResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.InsertOne(ctx, document, opts...)
		logCmd(wc.logMode, wc.cmd(c), "InsertOne", res, document)
		return err
	})
	return
}

func (wc *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.UpdateByID(ctx, id, update, opts...)
		logCmd(wc.logMode, wc.cmd(c), "UpdateByID", res, id, update)
		return err
	})
	return
}

func (wc *Collection) Name() string { return wc.coll.Name() }

func (wc *Collection) ReplaceOne(ctx context.Context, filter, replacement interface{}, opts ...*options.ReplaceOptions) (res *mongo.UpdateResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.ReplaceOne(ctx, filter, replacement, opts...)
		logCmd(wc.logMode, wc.cmd(c), "ReplaceOne", res, filter, replacement)
		return err
	})
	return
}

func (wc *Collection) UpdateMany(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.UpdateMany(ctx, filter, replacement, opts...)
		logCmd(wc.logMode, wc.cmd(c), "UpdateMany", res, filter, replacement)
		return err
	})
	return
}

func (wc *Collection) UpdateOne(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.UpdateOne(ctx, filter, replacement, opts...)
		logCmd(wc.logMode, wc.cmd(c), "UpdateOne", res, filter, replacement)
		return err
	})
	return
}

func (wc *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (res *mongo.ChangeStream, err error) {
	_ = wc.processor(func(c *cmd) error {
		res, err = wc.coll.Watch(ctx, pipeline, opts...)
		logCmd(wc.logMode, wc.cmd(c), "Watch", res, pipeline)
		return err
	})
	return
}

func (wc *Collection) Collection() *mongo.Collection {
	return wc.coll
}

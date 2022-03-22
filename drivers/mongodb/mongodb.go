package mongodb

import (
	"context"
	"fmt"
	"regexp"

	"github.com/k1LoW/tbls/dict"
	"github.com/k1LoW/tbls/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var re = regexp.MustCompile(`(?s)\n\s*`)

type Mongodb struct {
	ctx    context.Context
	client *mongo.Client
}

func New(ctx context.Context, client *mongo.Client) (*Mongodb, error) {
	return &Mongodb{
		ctx:    ctx,
		client: client,
	}, nil
}

func (d *Mongodb) Analyze(s *schema.Schema) error {
	drv, err := d.Info()
	if err != nil {
		return errors.WithStack(err)
	}
	s.Driver = drv

	tables := []*schema.Table{}
	dbValue := d.client.Database("test")
	colls, err := dbValue.ListCollectionSpecifications(d.ctx, bson.D{})
	if err != nil {
		return err
	}
	for _, coll := range colls {
		colVal := dbValue.Collection(coll.Name)
		indexes, err := d.listIndexes(colVal)
		if err != nil {
			return err
		}
		estimated, err := colVal.EstimatedDocumentCount(d.ctx)
		if err != nil {
			return err
		}
		table := &schema.Table{
			Name:    coll.Name,
			Type:    coll.Type,
			Indexes: indexes,
			Comment: fmt.Sprintf("Estimated count of document is %d", estimated),
		}
		tables = append(tables, table)
	}
	s.Tables = tables

	return nil
}

func (d *Mongodb) listIndexes(collection *mongo.Collection) ([]*schema.Index, error) {
	indexes := []*schema.Index{}
	indexSpec, err := collection.Indexes().ListSpecifications(d.ctx)
	if err != nil {
		return nil, err
	}
	for _, spec := range indexSpec {
		var isUnique string
		if spec.Unique != nil && *spec.Unique {
			isUnique = "Unique"
		} else {
			isUnique = "Non-unique"
		}
		comment := fmt.Sprintf("%s, Version %d", isUnique, spec.Version)
		Index := &schema.Index{
			Name:    spec.Name,
			Def:     spec.KeysDocument.String(),
			Comment: comment,
		}
		indexes = append(indexes, Index)
	}
	return indexes, nil
}

func (d *Mongodb) Info() (*schema.Driver, error) {
	dct := dict.New()
	dct.Merge(map[string]string{
		"Column":  "Attribute",
		"Columns": "Attributes",
		"Indexes": "Indexes",
	})

	driver := &schema.Driver{
		Name:            "mongodb",
		DatabaseVersion: "",
		Meta: &schema.DriverMeta{
			Dict: &dct,
		},
	}
	return driver, nil
}

package database

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/cfanatic/go-expenses/datasheet"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TIMEOUT = 3
)

type ttime = time.Time

type Database struct {
	addr string
	db   string
	opt  *options.ClientOptions
	clt  *mongo.Client
	coll *mongo.Collection
	Err  error
}

type Content struct {
	Date      ttime   `bson:"date"`
	Payee     string  `bson:"payee"`
	Desc      string  `bson:"desc"`
	Amount    float32 `bson:"amount"`
	Label     string  `bson:"label"`
	Hash      string  `bson:"hash"`
	Datasheet string  `bson:"datasheet"`
}

func New(address, database, collection string) *Database {
	db := &Database{addr: address, db: database}
	db.opt = options.Client().ApplyURI(address)
	if client, err := mongo.Connect(context.TODO(), db.opt); err == nil {
		db.clt = client
		db.coll = db.clt.Database(db.db).Collection(collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()
	if err := db.clt.Ping(ctx, nil); err == nil {
		db.Err = nil
	} else {
		db.Err = err
	}
	return db
}

func (db *Database) Document(key, value string) (Content, error) {
	doc := Content{}
	filt := bson.D{{Key: key, Value: value}}
	err := db.coll.FindOne(context.TODO(), filt).Decode(&doc)
	return doc, err
}

func (db *Database) Content(filter ...string) ([]Content, error) {
	docs := make([]Content, 0)
	cur := &mongo.Cursor{}
	query := func(filt interface{}) {
		cur, _ = db.coll.Find(context.TODO(), filt)
		defer cur.Close(context.TODO())
		for cur.Next(context.TODO()) {
			doc := Content{}
			cur.Decode(&doc)
			docs = append(docs, doc)
		}
	}
	switch {
	case len(filter) == 0:
		filt := bson.D{{}}
		query(filt)
	case len(filter) == 2:
		filt := bson.D{{Key: filter[0], Value: filter[1]}}
		query(filt)
	case len(filter) == 4:
		dateFrom, _ := time.Parse("01-02-06", filter[2])
		dateTo, _ := time.Parse("01-02-06", filter[3])
		filt := bson.M{}
		if len(filter[0]) == 0 && len(filter[1]) == 0 {
			filt = bson.M{"date": bson.M{"$gte": dateFrom, "$lte": dateTo}}
		} else {
			filt = bson.M{
				"$and": []interface{}{
					bson.D{{Key: filter[0], Value: filter[1]}},
					bson.M{"date": bson.M{"$gte": dateFrom, "$lte": dateTo}},
				}}
		}
		query(filt)
	}
	err := cur.Err()
	return docs, err
}

func (db *Database) Labels(label string) ([]interface{}, error) {
	values, err := db.coll.Distinct(context.TODO(), label, bson.D{{}})
	return values, err
}

func (db *Database) Save(document Content) error {
	hash := db.Hash(document)
	if _, err := db.Document("hash", hash); err != nil {
		document.Hash = hash
		_, err := db.coll.InsertOne(context.TODO(), document)
		return err
	} else {
		return nil
	}
}

func (db *Database) Update(old, new Content) error {
	doc := Content{}
	filt := bson.D{
		{Key: "date", Value: old.Date},
		{Key: "payee", Value: old.Payee},
		{Key: "desc", Value: old.Desc},
		{Key: "amount", Value: old.Amount},
	}
	if err := db.coll.FindOne(context.TODO(), filt).Decode(&doc); err == nil {
		db.coll.DeleteOne(context.TODO(), filt)
	}
	return db.Save(new)
}

func (db *Database) Delete() (int64, error) {
	res, err := db.coll.DeleteMany(context.TODO(), bson.D{{}})
	return res.DeletedCount, err
}

func (db *Database) Hash(content interface{}) string {
	var fp string

	switch data := content.(type) {
	case Content:
		fp = fmt.Sprintf("%s %s %s %f %s",
			data.Date,
			data.Payee,
			data.Desc,
			data.Amount,
			data.Label,
		)
	case datasheet.Content:
		fp = fmt.Sprintf("%s %s %s %f",
			data.Date,
			data.Payee,
			data.Desc,
			data.Amount,
		)
	}

	md5 := md5.Sum([]byte(fp))
	hash := fmt.Sprintf("%x", md5)
	return hash
}

func (db *Database) Print(content interface{}) {
	log.Printf("%+v\n", content)
}

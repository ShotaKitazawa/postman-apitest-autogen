package main

import (
	"encoding/json"
	"log"
	"os"

	myflag "github.com/ShotaKitazawa/postman-apitest-autogen/pkg/flag"
	"github.com/ShotaKitazawa/postman-apitest-autogen/pkg/model"
)

func main() {
	args := myflag.Parse()
	c, err := model.NewCollectionFromFile(args.CollectionFile)
	if err != nil {
		log.Fatal(err)
	}
	var f *model.Filters
	if args.FilterFile != "" {
		f, err = model.NewFilterFromFile(args.FilterFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	result := c
	for idx, item := range c.Item {
		if !item.IsOnlyOnceResponse() {
			log.Fatalf("in %s: responses are not only once", item.Name)
		}
		testset := model.NewTestset(item.Name)
		if err := testset.BindResponseStatus(item.Responses[0].Code); err != nil {
			log.Fatal(err)
		}
		if item.IsJsonResponse() {
			b := item.Responses[0].GetBody()
			if f != nil {
				b, err = f.Filter(item.Name, b)
				if err != nil {
					log.Fatal(err)
				}
			}
			if err := testset.BindResponseBody(b); err != nil {
				log.Fatal(err)
			}

		}
		result.Item[idx].Events = append([]model.CollectionEvent{}, model.NewCollectionEvent(testset.ToString()))
	}
	data, err := json.Marshal(&result)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(args.OutputFile, data, 0644); err != nil {
		log.Fatal(err)
	}
}

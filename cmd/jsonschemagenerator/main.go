package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Arthur1/dpr"
	"github.com/invopop/jsonschema"
)

func main() {
	schema := jsonschema.Reflect(new(dpr.Config))
	b, err := schema.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}

	// workaround: *jsonschema.Schema does not have MarshalIndent()
	var a any
	if err := json.Unmarshal(b, &a); err != nil {
		log.Fatalln(err)
	}
	pb, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	pb = append(pb, byte('\n'))

	if err := os.WriteFile("dprcconfig.schema.json", pb, 0666); err != nil {
		log.Fatalln(err)
	}
}

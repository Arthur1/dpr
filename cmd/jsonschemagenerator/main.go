package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Arthur1/dpr/internal/cli"
	"github.com/Arthur1/dpr/lifecyclepolicy"
	"github.com/invopop/jsonschema"
)

func main() {
	cases := map[any]string{
		new(cli.RawConfig):                   "dprcconfig.schema.json",
		new(lifecyclepolicy.LifecyclePolicy): "dprlifecyclepolicy.schema.json",
	}
	for obj, file := range cases {
		schema := jsonschema.Reflect(obj)
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

		if err := os.WriteFile(file, pb, 0666); err != nil {
			log.Fatalln(err)
		}
	}
}

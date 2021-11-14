package ldContext


import (
	"github.com/piprate/json-gold/ld"
	"sync"
	"fmt"
)


var ExpandOnce sync.Once
var CompactOnce sync.Once
var ldE *ld.RFC7324CachingDocumentLoader
var ldC *ld.RFC7324CachingDocumentLoader
var expand_lock sync.RWMutex

var compact_lock sync.RWMutex

//creating expand singleton object for document loader
func Expand_object() *ld.RFC7324CachingDocumentLoader {
        if ldE == nil {
                ExpandOnce.Do(
                        func() {
                                ldE = ld.NewRFC7324CachingDocumentLoader(nil)
                                _, err := ldE.LoadDocument("https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context-v1.3.jsonld")
                                fmt.Println("created object", ldE, err)
                        })
        } else {
                fmt.Println("The loader object is already created")
        }
        return ldE
}

//Expand 
func  ExpandEntity(v interface{}) ([]interface{}, error) {
        dl := Expand_object()
        proc := ld.NewJsonLdProcessor()
        opts := ld.NewJsonLdOptions("")
        opts.ProcessingMode = ld.JsonLd_1_1
        opts.DocumentLoader = dl
        expanded, err := proc.Expand(v, opts)
	fmt.Println("++++err++++++",err)
        return expanded, err
}




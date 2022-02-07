package main 

import "testing"


func test_matchLdMetadatas(t *testing.T) {
	ctxMeta := map[string]ContextMetadata {
			"location": 
		}
	actual := matchLdMetadatas()
	t.Errorf("Expected %v but got %v",  actual)
}

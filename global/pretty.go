//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package global

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//goland:noinspection GoUnusedExportedFunction
func JSONPretty(j []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, j, "", "\t")
	if err != nil {
		fmt.Printf("JSON parse error: %s\n", err.Error())
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}

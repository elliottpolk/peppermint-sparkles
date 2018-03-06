//  +build ignore

// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/uuid/gen_uuid.go
//

package main

import (
	"fmt"

	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/uuid"
)

func main() {
	fmt.Print(uuid.GetV4())
}

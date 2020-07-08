// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatacenter_GetTagsByCategories(t *testing.T) {
	ref := mor{Type: "type", Value: "val"}
	tags := []Tag{
		{
			Name:     "A",
			Category: "cat1",
		},
		{
			Name:     "B",
			Category: "cat1",
		},
		{
			Name:     "B",
			Category: "cat2",
		},
		{
			Name:     "A",
			Category: "cat2",
		},
	}
	tagsByObject := make(map[mor][]Tag)
	tagsByObject[ref] = tags
	dc := NewDatacenter(nil)
	dc.AddTags(tagsByObject)

	tbc := dc.GetTagsByCategories(ref)
	assert.Equal(t, "A|B", tbc["cat1"], "Tags should should be ordered")
	assert.Equal(t, "A|B", tbc["cat2"], "Tags should should be ordered")
}

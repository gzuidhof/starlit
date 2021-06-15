/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pages

import (
	"sort"

	"github.com/gohugoio/hugo/compare"
)

/* menuSorter adapted from Hugo's Menu sorter, originally Apache 2.0 licensed */

// A type to implement the sort interface for Menu
type menuSorter struct {
	menu Menu
	by   menuEntryBy
}


// Closure used in the Sort.Less method.
type menuEntryBy func(m1, m2 *MenuEntry) bool


func (by menuEntryBy) Sort(menu Menu) {
	ms := &menuSorter{
		menu: menu,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ms)
}

var defaultMenuEntrySort = func(m1, m2 *MenuEntry) bool {
	if m1.Weight == m2.Weight {
		c := compare.Strings(m1.Name, m2.Name)
		if c == 0 {
			return m1.Identifier < m2.Identifier
		}
		return c < 0
	}

	// if m2.Weight == 0 {
	// 	return true
	// }

	// if m1.Weight == 0 {
	// 	return false
	// }

	return m1.Weight < m2.Weight
}

func (ms *menuSorter) Len() int      { return len(ms.menu) }
func (ms *menuSorter) Swap(i, j int) { ms.menu[i], ms.menu[j] = ms.menu[j], ms.menu[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ms *menuSorter) Less(i, j int) bool { return ms.by(ms.menu[i], ms.menu[j]) }

// Sort sorts the menu by weight, name and then by identifier.
func (m Menu) Sort() Menu {
	menuEntryBy(defaultMenuEntrySort).Sort(m)
	return m
}

// Limit limits the returned menu to n entries.
func (m Menu) Limit(n int) Menu {
	if len(m) > n {
		return m[0:n]
	}
	return m
}

// ByWeight sorts the menu by the weight defined in the menu configuration.
func (m Menu) ByWeight() Menu {
	menuEntryBy(defaultMenuEntrySort).Sort(m)
	return m
}

// ByName sorts the menu by the name defined in the menu configuration.
func (m Menu) ByName() Menu {
	title := func(m1, m2 *MenuEntry) bool {
		return compare.LessStrings(m1.Name, m2.Name)
	}

	menuEntryBy(title).Sort(m)
	return m
}

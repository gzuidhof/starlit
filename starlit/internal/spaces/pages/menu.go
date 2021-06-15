/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pages

import (
	"strings"

	"github.com/gzuidhof/starlit/starlit/internal/content"
)

type MenuData struct {
	Entries Menu
	indexNodes map[string]*MenuEntry
	pathToEntry map[string]*MenuEntry
}

type Menu []*MenuEntry

 
type MenuEntry struct {
	Page *content.Page
	Identifier string
	Name string
	Href string
	title string
	Weight int

	Children Menu
}

func (m Menu) Add(me *MenuEntry) Menu {
	m = append(m, me)
	m.Sort()
	return m
}

// Reverse reverses the order of the menu entries.
func (m Menu) Reverse() Menu {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}

	return m
}

func (m *MenuEntry) Title() string {
	if m.title != "" {
		return m.title
	}

	if m.Page != nil {
		return m.Page.Title()
	}

	return ""
}

func (m *MenuData) ResolvePage(path string) *content.Page {
	pathWithoutFileExtension := strings.Split(path, ".")[0]
	me := m.pathToEntry[pathWithoutFileExtension]

	if me == nil {
		return nil
	}
	return me.Page
}
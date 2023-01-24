package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type TabsConfig map[string]fyne.CanvasObject

func CreateTabs(tabsConfig TabsConfig) *container.AppTabs {
	var tabs []*container.TabItem
	for name, content := range tabsConfig {
		tabs = append(tabs, container.NewTabItem(name, content))
	}
	tabsContainer := container.NewAppTabs(tabs...)
	tabsContainer.SetTabLocation(container.TabLocationTop)
	return tabsContainer
}

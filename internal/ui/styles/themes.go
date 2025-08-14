// Package styles defines color schemes and styles for the go-pwr UI.
package styles

import "github.com/charmbracelet/lipgloss"

// ColorScheme represents a color theme for the application.
type ColorScheme struct {
	Name      string
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Dim       lipgloss.Color
}

// Available color schemes
var (
	OceanBreeze = ColorScheme{
		Name:      "Ocean Breeze",
		Primary:   lipgloss.Color("39"),  // bright blue
		Secondary: lipgloss.Color("33"),  // blue
		Accent:    lipgloss.Color("45"),  // cyan
		Dim:       lipgloss.Color("244"), // gray
	}

	RocketPink = ColorScheme{
		Name:      "Rocket Pink",
		Primary:   lipgloss.Color("205"), // pink
		Secondary: lipgloss.Color("93"),  // purple
		Accent:    lipgloss.Color("198"), // bright pink
		Dim:       lipgloss.Color("244"), // gray
	}

	ForestNight = ColorScheme{
		Name:      "Forest Night",
		Primary:   lipgloss.Color("46"),  // green
		Secondary: lipgloss.Color("34"),  // forest green
		Accent:    lipgloss.Color("82"),  // lime
		Dim:       lipgloss.Color("244"), // gray
	}

	SunsetGlow = ColorScheme{
		Name:      "Sunset Glow",
		Primary:   lipgloss.Color("208"), // orange
		Secondary: lipgloss.Color("196"), // red
		Accent:    lipgloss.Color("226"), // yellow
		Dim:       lipgloss.Color("244"), // gray
	}

	PurpleHaze = ColorScheme{
		Name:      "Purple Haze",
		Primary:   lipgloss.Color("135"), // purple
		Secondary: lipgloss.Color("93"),  // violet
		Accent:    lipgloss.Color("171"), // magenta
		Dim:       lipgloss.Color("244"), // gray
	}

	ArcticFrost = ColorScheme{
		Name:      "Arctic Frost",
		Primary:   lipgloss.Color("51"),  // cyan
		Secondary: lipgloss.Color("39"),  // blue
		Accent:    lipgloss.Color("87"),  // light blue
		Dim:       lipgloss.Color("244"), // gray
	}
)

// AllSchemes returns all available color schemes.
func AllSchemes() []ColorScheme {
	return []ColorScheme{
		OceanBreeze,
		RocketPink,
		ForestNight,
		SunsetGlow,
		PurpleHaze,
		ArcticFrost,
	}
}

// GetSchemeByName returns a color scheme by name, or the default if not found.
func GetSchemeByName(name string) ColorScheme {
	for _, scheme := range AllSchemes() {
		if scheme.Name == name {
			return scheme
		}
	}
	// Return default scheme if not found
	return OceanBreeze
}

// Theme holds the current styling configuration.
type Theme struct {
	Current     ColorScheme
	Border      lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	TabBar      lipgloss.Style
}

// NewTheme creates a new theme with the given color scheme.
func NewTheme(scheme ColorScheme) *Theme {
	return &Theme{
		Current: scheme,
		Border: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2).
			BorderForeground(scheme.Primary),
		TabActive: lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Padding(0, 1).
			Foreground(scheme.Primary),
		TabInactive: lipgloss.NewStyle().
			Faint(true).
			Padding(0, 1).
			Foreground(scheme.Primary),
		TabBar: lipgloss.NewStyle().
			MarginBottom(1).
			Foreground(scheme.Primary),
	}
}

// UpdateScheme updates the theme with a new color scheme.
func (t *Theme) UpdateScheme(scheme ColorScheme) {
	t.Current = scheme
	t.Border = t.Border.BorderForeground(scheme.Primary)
	t.TabActive = t.TabActive.Foreground(scheme.Primary)
	t.TabInactive = t.TabInactive.Foreground(scheme.Primary)
	t.TabBar = t.TabBar.Foreground(scheme.Primary)
}

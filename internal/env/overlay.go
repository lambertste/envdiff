package env

import "fmt"

// OverlayOptions controls how layers are merged.
type OverlayOptions struct {
	// Overwrite allows later layers to overwrite keys from earlier layers.
	Overwrite bool
	// SkipEmpty skips keys with empty values from overlay layers.
	SkipEmpty bool
}

// DefaultOverlayOptions returns sensible defaults.
func DefaultOverlayOptions() OverlayOptions {
	return OverlayOptions{
		Overwrite: true,
		SkipEmpty: false,
	}
}

// Overlay merges a slice of Sets into a single Set, applying each layer in
// order. The first element is the base; subsequent elements are overlaid on top.
func Overlay(layers []*Set, opts OverlayOptions) (*Set, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("overlay: at least one layer required")
	}
	out := Clone(layers[0])
	for i, layer := range layers[1:] {
		if layer == nil {
			return nil, fmt.Errorf("overlay: layer %d is nil", i+1)
		}
		for _, k := range layer.Keys() {
			v, _ := layer.Get(k)
			if opts.SkipEmpty && v == "" {
				continue
			}
			_, exists := out.Get(k)
			if !exists || opts.Overwrite {
				out.Set(k, v)
			}
		}
	}
	return out, nil
}

// OverlayReport describes which keys were overwritten and which were preserved.
type OverlayReport struct {
	Overwritten []string
	Preserved   []string
	Added       []string
}

// OverlayWithReport performs an overlay and returns a report of changes.
func OverlayWithReport(base *Set, overlay *Set, opts OverlayOptions) (*Set, OverlayReport, error) {
	if base == nil {
		return nil, OverlayReport{}, fmt.Errorf("overlay: base is nil")
	}
	if overlay == nil {
		return nil, OverlayReport{}, fmt.Errorf("overlay: overlay is nil")
	}
	out := Clone(base)
	var report OverlayReport
	for _, k := range overlay.Keys() {
		v, _ := overlay.Get(k)
		if opts.SkipEmpty && v == "" {
			continue
		}
		_, exists := out.Get(k)
		if !exists {
			out.Set(k, v)
			report.Added = append(report.Added, k)
		} else if opts.Overwrite {
			out.Set(k, v)
			report.Overwritten = append(report.Overwritten, k)
		} else {
			report.Preserved = append(report.Preserved, k)
		}
	}
	return out, report, nil
}

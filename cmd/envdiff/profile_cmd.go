package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"envdiff/internal/profile"
)

const defaultProfilePath = ".envdiff/profiles.json"

func runProfileAdd(name, file string, tags []string) error {
	r, err := profile.Load(defaultProfilePath)
	if err != nil {
		return err
	}
	r.Add(profile.Profile{
		Name: name,
		File: file,
		Tags: tags,
	})
	if err := profile.Save(defaultProfilePath, r); err != nil {
		return err
	}
	fmt.Printf("profile %q added (file: %s)\n", name, file)
	return nil
}

func runProfileRemove(name string) error {
	r, err := profile.Load(defaultProfilePath)
	if err != nil {
		return err
	}
	if _, ok := r.Get(name); !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	r.Remove(name)
	if err := profile.Save(defaultProfilePath, r); err != nil {
		return err
	}
	fmt.Printf("profile %q removed\n", name)
	return nil
}

func runProfileList() error {
	r, err := profile.Load(defaultProfilePath)
	if err != nil {
		return err
	}
	names := r.List()
	if len(names) == 0 {
		fmt.Println("no profiles registered")
		return nil
	}
	sort.Strings(names)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tFILE\tTAGS")
	for _, n := range names {
		p, _ := r.Get(n)
		tags := "-"
		if len(p.Tags) > 0 {
			tags = fmt.Sprintf("%v", p.Tags)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.File, tags)
	}
	return w.Flush()
}

func runProfileShow(name string) error {
	r, err := profile.Load(defaultProfilePath)
	if err != nil {
		return err
	}
	p, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	fmt.Printf("Name: %s\nFile: %s\nTags: %v\n", p.Name, p.File, p.Tags)
	for k, v := range p.Meta {
		fmt.Printf("  %s: %s\n", k, v)
	}
	return nil
}

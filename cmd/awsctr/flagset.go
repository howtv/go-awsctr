package main

import "flag"

type Config struct {
	Region string
}

func NewFlagSet(c *Config) *flag.FlagSet {
	fs := flag.NewFlagSet("common", flag.ExitOnError)

	fs.StringVar(&c.Region, "region", defaultRegion, "aws region")

	return fs
}

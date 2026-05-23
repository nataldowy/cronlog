// Package config provides loading and validation of cronlog configuration
// files written in YAML.
//
// A minimal configuration only needs to specify a store path; all other
// fields have sensible defaults:
//
//	store:
//	  path: /var/log/cronlog
//	notify:
//	  smtp_host: localhost
//	  smtp_port: 25
//	  from: cronlog@example.com
//	  recipients:
//	    - ops@example.com
//	runner:
//	  default_timeout: 24h
//
// Load returns a validated *Config or an error describing the first
// problem encountered.
package config

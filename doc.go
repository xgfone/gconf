// Package gconf is an extensible and powerful go configuration manager,
// which is inspired by https://github.com/openstack/oslo.config.
//
// Goal
//
// This package is aimed at
//
//   1. A atomic key-value configuration center with the multi-group and the option.
//   2. Support the multi-parser to parse the configurations from many sources
//      with the different format.
//   3. Change the configuration dynamically during running and watch it.
//   4. Observe the change of the configuration.
//
package gconf

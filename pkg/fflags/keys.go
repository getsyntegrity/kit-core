package fflags

// Canonical flag key constants used across multiple services are defined in this file.
// Use these constants instead of string literals when a key is shared by BFF and control-plane services.
// Do not add provider SDK imports or service-specific keys; those belong in each service's internal/fflags.
//
// (No shared keys are currently defined here; service-owned keys live in internal/fflags/keys.go per service.)

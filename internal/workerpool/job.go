// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package workerpool just a wrapper to faciliate the workerpool struct
package workerpool

type Result struct {
	Value any
	Err   error
}

type Job func() Result

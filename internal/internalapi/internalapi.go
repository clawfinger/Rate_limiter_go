package internalapi

import "context"

type OperationResult struct {
	Ok     bool
	Reason string
}

type StorageResult struct {
	Status Restriction
	Reason string
	Err    error
}

type Restriction uint8

const (
	NotSet      Restriction = 0
	Whitelisted Restriction = 1
	Blacklisted Restriction = 2
)

type AbstractRateManager interface {
	Manage(ip string, login string, password string) *OperationResult
	DropIPStats(ip string)
	DropLiginStats(login string)
	DropPasswordStats(password string)
}

type AbstractStorage interface {
	CheckIP(ctx context.Context, ip string) *StorageResult
	SetIP(ctx context.Context, ip string, restriction Restriction) error
	RemoveIP(ctx context.Context, ip string, restriction Restriction) error
}

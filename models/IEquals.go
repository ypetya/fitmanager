package models

type IEquals interface {
	equals(o IEquals) bool
}

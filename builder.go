package gsd

type builder interface {
	BuildSelect(ctx *buildContext, info *selectInfo) error
	BuildInsert(ctx *buildContext, info *insertInfo) error
	BuildUpdate(ctx *buildContext, info *updateInfo) error
	BuildDelete(ctx *buildContext, info *deleteInfo) error
}

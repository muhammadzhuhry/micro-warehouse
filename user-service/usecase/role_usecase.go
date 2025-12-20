package usecase

import (
	"context"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/repository"
)

type RoleUsecaseInterface interface {
	CreateRole(ctx context.Context, role model.Role) error
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type RoleUsecase struct {
	roleRepo repository.RoleRepositoryInterface
}

func NewRoleUsecase(roleRepo repository.RoleRepositoryInterface) RoleUsecaseInterface {
	return &RoleUsecase{
		roleRepo: roleRepo,
	}
}

// CreateRole implements RoleUsecaseInterface.
func (r *RoleUsecase) CreateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.CreateRole(ctx, role)
}

// DeleteRole implements RoleUsecaseInterface.
func (r *RoleUsecase) DeleteRole(ctx context.Context, id uint) error {
	return r.roleRepo.DeleteRole(ctx, id)
}

// GetAllRoles implements RoleUsecaseInterface.
func (r *RoleUsecase) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetAllRoles(ctx)
}

// GetRoleByID implements RoleUsecaseInterface.
func (r *RoleUsecase) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

// UpdateRole implements RoleUsecaseInterface.
func (r *RoleUsecase) UpdateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.UpdateRole(ctx, role)
}

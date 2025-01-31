import { Permission, Role, RolePermissionEntity } from '@crm/pages/staff/pages/all/crm-staff.model';

export interface RoleAndPermissionPayload {
  user_id: string
  role_permission: RolePermissionEntity[];
}

export interface DeleteRole {
  user_id: string
  role: Role
}

export interface DeletePermission {
  user_id: string
  role: Role
  permission: Permission
}

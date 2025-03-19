import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import {
  CRMStaffModel,
  Permission,
  Role,
  RolePermissionEntity
} from '@crm/pages/staff/pages/all/crm-staff.model';
import { Divider } from 'primeng/divider';
import { Button } from 'primeng/button';
import { Dialog } from 'primeng/dialog';
import { NewRolePermissionComponent } from './ui/new-role-permission.component';
import {
  DeletePermission,
  DeleteRole,
  RoleAndPermissionPayload
} from '@crm/pages/account/account.model';
import { ApiState } from '@root/app.model';
import { NewPermissionComponent } from './ui/new-permission.component';
import { ConfirmationService } from 'primeng/api';
import { ConfirmPopup } from 'primeng/confirmpopup';

@Component({
  selector: 'app-role',
  imports: [
    Divider,
    Button,
    Dialog,
    NewRolePermissionComponent,
    NewPermissionComponent,
    ConfirmPopup
  ],
  providers: [ConfirmationService],
  templateUrl: './role.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoleComponent {
  private readonly confirmService = inject(ConfirmationService);

  staff = input.required<CRMStaffModel>();
  createRolesAndPermissionState = input.required<ApiState>();
  addPermissionsState = input.required<ApiState>();
  deleteRoleState = input.required<ApiState>();
  deletePermissionState = input.required<ApiState>();

  readonly createRolesAndPermissionsEmitter =
    output<RoleAndPermissionPayload>();
  readonly addPermissionsEmitter = output<RoleAndPermissionPayload>();
  readonly deleteRoleEmitter = output<DeleteRole>();
  readonly deletePermissionEmitter = output<DeletePermission>();

  protected readonly dropdown: boolean[] = [];
  protected toggleNewRolePermission = false;
  protected toggleInsertPermission = false;
  protected selectRoleAndPermissions: RolePermissionEntity | undefined;

  protected readonly filter = (a: RolePermissionEntity[]) =>
    [Role.STAFF, Role.DEVELOPER].filter(r => !a.some(rp => rp.role === r));

  protected readonly permissions = (srp: RolePermissionEntity) =>
    [Permission.READ, Permission.WRITE, Permission.DELETE].filter(
      p => !srp.permissions.some(rp => rp === p)
    );

  protected readonly confirmDeleteRole = (
    event: Event,
    userId: string,
    role: Role
  ) =>
    this.confirmService.confirm({
      target: event.target as EventTarget,
      message: 'Are you sure you want to unlink the role?',
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptButtonProps: { label: 'Unlink' },
      accept: () => this.deleteRoleEmitter.emit({ user_id: userId, role: role })
    });

  protected readonly confirmDeletePermission = (
    event: Event,
    userId: string,
    role: Role,
    perm: Permission
  ) =>
    this.confirmService.confirm({
      target: event.target as EventTarget,
      message: `Are you sure you want to unlink the permission?`,
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptButtonProps: { label: 'Unlink' },
      accept: () =>
        this.deletePermissionEmitter.emit({
          user_id: userId,
          role: role,
          permission: perm
        })
    });
}

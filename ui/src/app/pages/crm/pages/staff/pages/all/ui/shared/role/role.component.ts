import {
  ChangeDetectionStrategy,
  Component,
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
import { RoleAndPermissionPayload } from '@crm/pages/account/account.model';
import { ApiState } from '@root/app.model';
import { NewPermissionComponent } from './ui/new-permission.component';

@Component({
  selector: 'app-role',
  imports: [
    Divider,
    Button,
    Dialog,
    NewRolePermissionComponent,
    NewPermissionComponent
  ],
  templateUrl: './role.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoleComponent {
  staff = input.required<CRMStaffModel>();
  createRolesAndPermissionState = input.required<ApiState>();

  readonly createRolesAndPermissionsEmitter =
    output<RoleAndPermissionPayload>();

  protected readonly state = ApiState;
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
}

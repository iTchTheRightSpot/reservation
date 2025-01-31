import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Divider } from 'primeng/divider';
import { Button } from 'primeng/button';
import { Dialog } from 'primeng/dialog';
import { NewRolePermissionComponent } from './ui/new-role-permission.component';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-role',
  imports: [Divider, Button, Dialog, NewRolePermissionComponent, FormsModule],
  templateUrl: './role.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoleComponent {
  staff = input.required<CRMStaffModel>();

  protected readonly dropdown: boolean[] = [];
  protected toggleNewRolePermission = false;
}

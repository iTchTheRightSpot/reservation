import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import {
  CRMStaffModel,
  Role,
  RolePermissionEntity
} from '@crm/pages/staff/pages/all/crm-staff.model';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { Avatar } from 'primeng/avatar';
import { Textarea } from 'primeng/textarea';

@Component({
  selector: 'app-view-staff',
  imports: [
    IconField,
    InputIcon,
    InputText,
    Select,
    FormsModule,
    Avatar,
    Textarea
  ],
  templateUrl: './view-staff.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ViewStaffComponent {
  staff = input.required<CRMStaffModel>();

  protected selectedRole: Role | undefined;

  protected readonly roles = (ac: RolePermissionEntity[]) =>
    ac.map(r => r.role);

  protected readonly permissions = (ac: RolePermissionEntity[]) => {
    const r = this.selectedRole;
    if (!r) return [];
    return ac
      .filter(a => a.role === r)
      .flatMap(rp => rp.permissions.map(p => p));
  };
}

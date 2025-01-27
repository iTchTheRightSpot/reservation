import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';

@Component({
  selector: 'app-role',
  imports: [],
  templateUrl: './role.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RoleComponent {
  staff = input.required<CRMStaffModel>();
}

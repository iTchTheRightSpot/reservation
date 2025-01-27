import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/crm-staff.model';

@Component({
  selector: 'app-staff-detail-holder',
  imports: [],
  template: `all staffs details component {{ staff().firstname }}`,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffDetailComponent {
  staff = input.required<CRMStaffModel>();
}

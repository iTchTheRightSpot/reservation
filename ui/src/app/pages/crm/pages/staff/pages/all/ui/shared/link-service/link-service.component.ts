import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';

@Component({
  selector: 'app-link-service',
  imports: [],
  templateUrl: './link-service.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LinkServiceComponent {
  staff = input.required<CRMStaffModel>();
}

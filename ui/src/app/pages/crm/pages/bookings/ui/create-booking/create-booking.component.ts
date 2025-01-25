import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/crm-staff.model';

@Component({
  selector: 'app-create-booking',
  imports: [],
  templateUrl: './create-booking.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CreateBookingComponent {
  staffs = input.required<CRMStaffModel[]>()
  protected readonly services = []
}

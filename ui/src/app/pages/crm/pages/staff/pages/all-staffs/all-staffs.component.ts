import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { CRMStaffsService } from '@crm/pages/staff/crm-staffs.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { CRMStaffModel } from '@crm/pages/staff/crm-staff.model';

@Component({
  selector: 'app-all-staffs',
  imports: [],
  templateUrl: './all-staffs.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AllStaffsComponent {
  private readonly service = inject(CRMStaffsService);

  protected readonly staffs = toSignal(this.service.staffs(), {
    initialValue: { state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>
  });
}

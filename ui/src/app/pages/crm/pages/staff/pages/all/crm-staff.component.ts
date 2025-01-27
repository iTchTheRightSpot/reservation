import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { CRMStaffsService } from './crm-staff.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { CRMStaffModel } from './crm-staff.model';
import { Button } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { Avatar } from 'primeng/avatar';
import { Tag } from 'primeng/tag';
import { Router } from '@angular/router';
import { CRM_ROUTE } from '@root/app.routes';
import { CRM_STAFFS_ROUTE } from '@crm/crm.routes';
import { REGISTER_ROUTE } from '@crm/pages/staff/staff-core.routes';
import { Drawer } from 'primeng/drawer';
import { StaffDetailComponent } from './ui/staff-detail.component';

@Component({
  selector: 'app-crm-staff',
  imports: [Button, TableModule, Avatar, Tag, Drawer, StaffDetailComponent],
  templateUrl: './crm-staff.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CrmStaffComponent {
  private readonly service = inject(CRMStaffsService);
  private readonly router = inject(Router);

  protected first = 0;
  protected rows = 5;
  protected readonly state = ApiState;
  protected toggleStaffDetails = false;
  protected readonly thead = ['Image', 'Firstname', 'Lastname', 'Account'];
  protected selectedStaff: CRMStaffModel | undefined;

  protected readonly register = async () =>
    await this.router.navigate([
      `${CRM_ROUTE}/${CRM_STAFFS_ROUTE}/${REGISTER_ROUTE}`
    ]);

  protected readonly staffs = toSignal(this.service.staffs(), {
    initialValue: { state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>
  });
}

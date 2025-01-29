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
import { map, Subject, switchMap, tap } from 'rxjs';
import { ServiceTypeToStaffPayload } from './ui/shared/link-service/link-service.model';
import { CRMServiceTypeService } from '@crm/pages/service-type/crm-service-type.service';
import { ScheduleService } from '@crm/pages/account/pages/schedule/schedule.service';
import { StaffScheduleEmitter } from './ui/shared/staff-schedule/staff-schedule.model';
import { Schedule } from '@crm/pages/account/pages/schedule/schedule.model';
import { CreateUpdateScheduleModel } from '@crm/pages/staff/pages/all/ui/shared/staff-schedule/ui/shared/shared-create-update-schedule.model';

@Component({
  selector: 'app-crm-staff',
  imports: [Button, TableModule, Avatar, Tag, Drawer, StaffDetailComponent],
  templateUrl: './crm-staff.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CrmStaffComponent {
  private readonly staffsService = inject(CRMStaffsService);
  private readonly serviceType = inject(CRMServiceTypeService);
  private readonly scheduleService = inject(ScheduleService);
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

  protected readonly allServices = toSignal(
    this.serviceType.all().pipe(
      map(
        o =>
          ({
            state: o.state,
            message: o.message,
            data: o.data?.map(s => s.name)
          }) as ApiResponse<string[]>
      )
    ),
    { initialValue: { state: ApiState.LOADING } as ApiResponse<string[]> }
  );

  protected readonly staffs = toSignal(this.staffsService.staffs(), {
    initialValue: { state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>
  });

  protected readonly linkServiceToStaffEmitter =
    new Subject<ServiceTypeToStaffPayload>();
  protected readonly serviceToStaffState = toSignal(
    this.linkServiceToStaffEmitter.asObservable().pipe(
      switchMap(o =>
        this.serviceType.linkServiceToStaff(o).pipe(
          tap(() => {
            CRMServiceTypeService.ServiceTypesByStaffCache.clear();
            this.servicesByStaffEmitter.next(o.staff_id);
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly servicesByStaffEmitter = new Subject<string>();
  protected readonly servicesByStaff = toSignal(
    this.servicesByStaffEmitter
      .asObservable()
      .pipe(switchMap(o => this.serviceType.servicesByStaff(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<string[]> }
  );

  protected readonly deLinkServiceFromStaffEmitter =
    new Subject<ServiceTypeToStaffPayload>();
  protected readonly deLinkServiceFromStaffState = toSignal(
    this.deLinkServiceFromStaffEmitter
      .asObservable()
      .pipe(
        switchMap(o =>
          this.serviceType
            .deLinkServiceFromStaff(o)
            .pipe(tap(() => this.servicesByStaffEmitter.next(o.staff_id)))
        )
      ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly staffSchedulesEmitter =
    new Subject<StaffScheduleEmitter>();
  protected readonly staffSchedules = toSignal(
    this.staffSchedulesEmitter
      .asObservable()
      .pipe(switchMap(o => this.scheduleService.schedulesByStaff(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<Schedule[]> }
  );

  protected readonly createScheduleEmitter =
    new Subject<CreateUpdateScheduleModel>();
  protected readonly createSchedule = toSignal(
    this.createScheduleEmitter
      .asObservable()
      .pipe(switchMap(o => this.scheduleService.create(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly updateScheduleEmitter =
    new Subject<CreateUpdateScheduleModel>();
  protected readonly updateSchedule = toSignal(
    this.createScheduleEmitter
      .asObservable()
      .pipe(switchMap(o => this.scheduleService.update(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly deleteScheduleEmitter = new Subject<number>();
  protected readonly deleteSchedule = toSignal(
    this.deleteScheduleEmitter
      .asObservable()
      .pipe(switchMap(o => this.scheduleService.delete(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}

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
import { RootRoutes } from '@root/app.routes';
import { CRMRoutes } from '@crm/crm.routes';
import { StaffRoutes } from '@crm/pages/staff/staff-core.routes';
import { Drawer } from 'primeng/drawer';
import { StaffDetailComponent } from './ui/staff-detail.component';
import { map, Subject, switchMap, tap } from 'rxjs';
import { ServiceTypeToStaffPayload } from './ui/shared/link-service/link-service.model';
import { CRMServiceTypeService } from '@crm/pages/service-type/crm-service-type.service';
import { ScheduleService } from '@crm/pages/account/pages/schedule/schedule.service';
import {
  DeleteScheduleModel,
  StaffScheduleEmitter
} from './ui/shared/staff-schedule/staff-schedule.model';
import {
  CreateScheduleModel,
  Schedule,
  UpdateScheduleModel
} from '@crm/pages/account/pages/schedule/schedule.model';
import { AccountService } from '@crm/pages/account/account.service';
import {
  DeletePermission,
  DeleteRole,
  RoleAndPermissionPayload
} from '@crm/pages/account/account.model';

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
  private readonly accountService = inject(AccountService);

  protected first = 0;
  protected rows = 5;
  protected readonly state = ApiState;
  protected toggleStaffDetails = false;
  protected readonly thead = ['Image', 'Firstname', 'Lastname', 'Account'];
  protected selectedStaff: CRMStaffModel | undefined;

  protected readonly register = async () =>
    await this.router.navigate([
      `${RootRoutes.CRM}/${CRMRoutes.STAFFS}/${StaffRoutes.REGISTER}`
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

  protected readonly createScheduleEmitter = new Subject<CreateScheduleModel>();
  protected readonly createSchedule = toSignal(
    this.createScheduleEmitter.asObservable().pipe(
      switchMap(o =>
        this.scheduleService.create(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.staff_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly updateScheduleEmitter = new Subject<UpdateScheduleModel>();
  protected readonly updateSchedule = toSignal(
    this.updateScheduleEmitter.asObservable().pipe(
      switchMap(o =>
        this.scheduleService.update(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.staff_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly deleteScheduleEmitter = new Subject<DeleteScheduleModel>();
  protected readonly deleteSchedule = toSignal(
    this.deleteScheduleEmitter.asObservable().pipe(
      switchMap(o =>
        this.scheduleService.delete(o.schedule_id).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.staff_id,
                date: o.date,
                page: o.page,
                size: o.size
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly createRolesAndPermissionsEmitter =
    new Subject<RoleAndPermissionPayload>();
  protected readonly createRolesAndPermissions = toSignal(
    this.createRolesAndPermissionsEmitter.asObservable().pipe(
      switchMap(o =>
        this.accountService.createRoleAndPermission(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.user_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly addPermissionsEmitter =
    new Subject<RoleAndPermissionPayload>();
  protected readonly addPermissionsState = toSignal(
    this.addPermissionsEmitter.asObservable().pipe(
      switchMap(o =>
        this.accountService.createRoleAndPermission(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.user_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly deleteRoleEmitter = new Subject<DeleteRole>();
  protected readonly deleteRoleState = toSignal(
    this.deleteRoleEmitter.asObservable().pipe(
      switchMap(o =>
        this.accountService.deleteRole(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.user_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly deletePermissionEmitter = new Subject<DeletePermission>();
  protected readonly deletePermissionState = toSignal(
    this.deletePermissionEmitter.asObservable().pipe(
      switchMap(o =>
        this.accountService.deletePermission(o).pipe(
          tap(s => {
            if (s.state === ApiState.LOADED) {
              ScheduleService.SchedulesByStaffCache.clear();
              ScheduleService.AllSchedulesCache.clear();
              this.staffSchedulesEmitter.next({
                staff_id: o.user_id,
                date: new Date(),
                page: 0,
                size: 10
              });
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}

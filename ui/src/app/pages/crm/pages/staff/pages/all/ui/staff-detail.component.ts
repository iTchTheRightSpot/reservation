import {
  ChangeDetectionStrategy,
  Component,
  input,
  output
} from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Tab, TabList, TabPanel, TabPanels, Tabs } from 'primeng/tabs';
import { ViewStaffComponent } from './shared/view/view-staff.component';
import { Tag } from 'primeng/tag';
import { RoleComponent } from './shared/role/role.component';
import { LinkServiceComponent } from './shared/link-service/link-service.component';
import { ServiceTypeToStaffPayload } from './shared/link-service/link-service.model';
import { ApiResponse, ApiState } from '@root/app.model';
import { StaffScheduleComponent } from './shared/staff-schedule/staff-schedule.component';
import {
  DeleteScheduleModel,
  StaffScheduleEmitter
} from '@crm/pages/staff/pages/all/ui/shared/staff-schedule/staff-schedule.model';
import {
  CreateScheduleModel,
  Schedule,
  UpdateScheduleModel
} from '@crm/pages/account/pages/schedule/schedule.model';

@Component({
  selector: 'app-staff-detail-holder',
  imports: [
    Tabs,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
    ViewStaffComponent,
    Tag,
    RoleComponent,
    LinkServiceComponent,
    StaffScheduleComponent
  ],
  template: `
    <p-tabs value="0">
      <p-tablist>
        <p-tab value="0">
          Account
          @if (staff().locked) {
            <p-tag [value]="'Locked'" [severity]="'danger'" />
          }
        </p-tab>
        <p-tab value="1">Roles & Permissions</p-tab>
        <p-tab value="2">Services</p-tab>
        <p-tab value="3">Schedules</p-tab>
      </p-tablist>
      <p-tabpanels>
        <p-tabpanel value="0">
          <app-view-staff [staff]="staff()" />
        </p-tabpanel>
        <p-tabpanel value="1">
          <app-role [staff]="staff()" />
        </p-tabpanel>
        <p-tabpanel value="2">
          <app-link-service
            [allServices]="allServices()"
            [linkServiceLoadingState]="linkServiceToStaffLoadingState()"
            [staffId]="staff().user_id"
            [staffServices]="servicesByStaff()"
            [deLinkServiceFromStaffState]="deLinkServiceFromStaffState()"
            (deLinkServiceFromStaffEmitter)="
              deLinkServiceFromStaffEmitter.emit($event)
            "
            (linkServiceToStaffEmitter)="linkServiceToStaffEmitter.emit($event)"
          />
        </p-tabpanel>
        <p-tabpanel value="3">
          <app-staff-schedule
            [staffId]="staff().user_id"
            [schedules]="staffSchedules()"
            [createScheduleLoadingState]="createScheduleLoadingState()"
            [updateScheduleLoadingState]="updateScheduleLoadingState()"
            (dateClicked)="schedulesEmitter.emit($event)"
            (createScheduleEmitter)="createScheduleEmitter.emit($event)"
            (updateScheduleEmitter)="updateScheduleEmitter.emit($event)"
            (deleteScheduleEmitter)="deleteScheduleEmitter.emit($event)"
          />
        </p-tabpanel>
      </p-tabpanels>
    </p-tabs>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffDetailComponent {
  staff = input.required<CRMStaffModel>();
  linkServiceToStaffLoadingState = input.required<ApiState>();
  servicesByStaff = input.required<ApiResponse<string[]>>();
  allServices = input.required<ApiResponse<string[]>>();
  deLinkServiceFromStaffState = input.required<ApiState>();
  staffSchedules = input.required<ApiResponse<Schedule[]>>();
  createScheduleLoadingState = input.required<ApiState>();
  updateScheduleLoadingState = input.required<ApiState>();
  deleteScheduleState = input.required<ApiState>();

  readonly deLinkServiceFromStaffEmitter = output<ServiceTypeToStaffPayload>();
  readonly linkServiceToStaffEmitter = output<ServiceTypeToStaffPayload>();
  readonly schedulesEmitter = output<StaffScheduleEmitter>();
  readonly createScheduleEmitter = output<CreateScheduleModel>();
  readonly updateScheduleEmitter = output<UpdateScheduleModel>();
  readonly deleteScheduleEmitter = output<DeleteScheduleModel>();
}

import {
  ChangeDetectionStrategy,
  Component,
  input,
  output
} from '@angular/core';
import { Badge } from 'primeng/badge';
import { Button } from 'primeng/button';
import { DatePicker } from 'primeng/datepicker';
import { FloatLabel } from 'primeng/floatlabel';
import { TableModule } from 'primeng/table';
import { ApiResponse, ApiState } from '@root/app.model';
import { Schedule } from '@crm/pages/account/pages/schedule/schedule.model';
import { StaffScheduleEmitter } from './staff-schedule.model';
import { Dialog } from 'primeng/dialog';
import { CreateScheduleComponent } from './ui/create-schedule.component';
import { EditScheduleComponent } from './ui/edit-schedule.component';
import { CreateUpdateScheduleModel } from './ui/shared/shared-create-update-schedule.model';

@Component({
  selector: 'app-staff-schedule',
  imports: [
    Badge,
    Button,
    DatePicker,
    FloatLabel,
    TableModule,
    Dialog,
    CreateScheduleComponent,
    EditScheduleComponent
  ],
  templateUrl: './staff-schedule.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffScheduleComponent {
  staffId = input.required<string>();
  schedules = input.required<ApiResponse<Schedule[]>>();
  createScheduleLoadingState = input.required<ApiState>();
  updateScheduleLoadingState = input.required<ApiState>();

  readonly dateClicked = output<StaffScheduleEmitter>();
  readonly createScheduleEmitter = output<CreateUpdateScheduleModel>();
  readonly updateScheduleEmitter = output<CreateUpdateScheduleModel>();
  readonly deleteScheduleEmitter = output<number>();

  protected first = 0;
  protected rows = 5;
  protected date = new Date();
  protected readonly state = ApiState;
  protected toggleCreateSchedule = false;
  protected toggleEditSchedule = false;
  protected readonly thead = [
    'From',
    'To',
    'Visible (show to customers)',
    'Reoccurring (weekly)'
  ];

  protected readonly fm = (d: number) =>
    new Date(d).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour12: true,
      hour: 'numeric',
      minute: 'numeric'
    });
}

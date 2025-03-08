import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { Badge } from 'primeng/badge';
import { Button } from 'primeng/button';
import { DatePicker } from 'primeng/datepicker';
import { FloatLabel } from 'primeng/floatlabel';
import { TableModule } from 'primeng/table';
import { ApiResponse, ApiState } from '@root/app.model';
import {
  CreateScheduleModel,
  Schedule,
  UpdateScheduleModel
} from '@crm/pages/account/pages/schedule/schedule.model';
import {
  DeleteScheduleModel,
  StaffScheduleEmitter
} from './staff-schedule.model';
import { Dialog } from 'primeng/dialog';
import { CreateScheduleComponent } from './ui/create-schedule.component';
import { EditScheduleComponent } from './ui/edit-schedule.component';
import { ConfirmPopup } from 'primeng/confirmpopup';
import { ConfirmationService } from 'primeng/api';

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
    EditScheduleComponent,
    ConfirmPopup
  ],
  providers: [ConfirmationService],
  templateUrl: './staff-schedule.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffScheduleComponent {
  protected readonly service = inject(ConfirmationService);

  staffId = input.required<string>();
  schedules = input.required<ApiResponse<Schedule[]>>();
  createScheduleLoadingState = input.required<ApiState>();
  updateScheduleLoadingState = input.required<ApiState>();

  readonly dateClicked = output<StaffScheduleEmitter>();
  readonly createScheduleEmitter = output<CreateScheduleModel>();
  readonly updateScheduleEmitter = output<UpdateScheduleModel>();
  readonly deleteScheduleEmitter = output<DeleteScheduleModel>();

  protected first = 0;
  protected rows = 5;
  protected date = new Date();
  protected readonly state = ApiState;
  protected selectedSchedule: Schedule | undefined;
  protected toggleCreateSchedule = false;
  protected toggleEditSchedule = false;
  protected readonly thead = [
    'From',
    'To',
    'Visible (show to customers)',
    'Reoccurring (weekly)'
  ];

  protected readonly fm = (d: number) =>
    new Date(Number(d)).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour12: true,
      hour: 'numeric',
      minute: 'numeric'
    });

  protected readonly delete = (
    event: Event,
    scheduleId: number,
    staffId: string
  ) => {
    this.service.confirm({
      target: event.target as EventTarget,
      message: 'Are you sure you want to complete deletion?',
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptButtonProps: { label: 'Delete' },
      accept: () =>
        this.deleteScheduleEmitter.emit({
          schedule_id: scheduleId,
          staff_id: staffId,
          page: this.first,
          size: this.rows,
          date: this.date
        })
    });
  };
}

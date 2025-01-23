import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { TableModule } from 'primeng/table';
import { FloatLabel } from 'primeng/floatlabel';
import { DatePicker } from 'primeng/datepicker';
import { Button } from 'primeng/button';
import { ScheduleService } from './schedule.service';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  BehaviorSubject,
  debounceTime,
  distinctUntilChanged,
  switchMap
} from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { Schedule } from './schedule.model';
import { tabledata } from '@crm/shared/data-access/crm.util';
import { Badge } from 'primeng/badge';

@Component({
  selector: 'app-schedule',
  imports: [TableModule, FloatLabel, DatePicker, Button, Badge],
  templateUrl: './schedule.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ScheduleComponent {
  private readonly service = inject(ScheduleService);

  protected first = 0;
  protected rows = 10;
  protected date = new Date();
  protected readonly state = ApiState;
  protected readonly thead = [
    'From',
    'To',
    'Visible (show to customers)',
    'Reoccurring (weekly)'
  ];

  protected readonly emitter = new BehaviorSubject<{
    date: Date;
    page: number;
    size: number;
  }>({ date: new Date(), page: this.first, size: this.rows });

  protected readonly models = toSignal(
    this.emitter.asObservable().pipe(
      distinctUntilChanged(
        (prev, curr) =>
          prev.date.toDateString() === curr.date.toDateString() &&
          prev.page === curr.page &&
          prev.size === curr.size
      ),
      debounceTime(700),
      switchMap(o => this.service.all(o.date, o.page, o.size))
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<Schedule[]> }
  );

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

  protected readonly data = (m: ApiResponse<Schedule[]>) =>
    tabledata<Schedule[]>(m);
}

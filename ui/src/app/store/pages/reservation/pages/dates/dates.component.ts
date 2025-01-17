import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal
} from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { DatesService } from './dates.service';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  BehaviorSubject,
  debounceTime,
  distinctUntilChanged,
  switchMap,
  tap
} from 'rxjs';
import { ApiResponse, ApiState, TIMEZONE } from '@root/app.util';
import { DateModel } from '@store/pages/reservation/pages/dates/dates.model';
import { FormsModule } from '@angular/forms';
import { DatePicker, DatePickerYearChangeEvent } from 'primeng/datepicker';
import { SummaryComponent } from '@store/pages/reservation/shared/summary.component';
import { ReservationModel } from '@store/pages/reservation/reservation.model';
import { Skeleton } from 'primeng/skeleton';
import { Message } from 'primeng/message';
import moment from 'moment-timezone';

@Component({
  selector: 'app-dates',
  imports: [FormsModule, DatePicker, SummaryComponent, Skeleton, Message],
  templateUrl: './dates.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DatesComponent {
  protected readonly service1 = inject(ReservationService);
  private readonly service2 = inject(DatesService);

  protected readonly state = ApiState;
  protected readonly today = new Date();
  protected date: Date[] | undefined;
  protected readonly tz = TIMEZONE;

  protected readonly selected = new BehaviorSubject(new Date());

  protected readonly style = (len: number) =>
    Array.from({ length: len }, (_, index) => index);

  protected readonly invalidDates = signal<Date[]>([]);
  protected readonly req = toSignal(
    this.selected.asObservable().pipe(
      distinctUntilChanged(),
      debounceTime(800),
      switchMap(date =>
        this.service2.dates(date).pipe(
          tap(obj => {
            if (obj.state === ApiState.LOADED && obj.data) {
              this.invalidDates.set(this.filter(date, obj.data));
            }
          })
        )
      )
    ),
    {
      initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
        DateModel[]
      >
    }
  );

  protected readonly timeformat: Intl.DateTimeFormatOptions = {
    weekday: 'long',
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  };

  protected readonly check = (o: ReservationModel) =>
    !o.services || o.services.length < 1 || !o.staff
      ? undefined
      : { services: o.services, staff: o.staff };

  private readonly filter = (date: Date, valid: DateModel[]) => {
    const validDates = valid.map(d =>
      moment.tz(Number(d.date), TIMEZONE).toDate()
    );

    const endOfMonth = new Date(date.getFullYear(), date.getMonth() + 1, 0);
    const daysInMonth = endOfMonth.getDate();

    const allDatesInMonth = Array.from(
      { length: daysInMonth },
      (_, index) => new Date(date.getFullYear(), date.getMonth(), index + 1)
    );

    return allDatesInMonth.filter(
      date =>
        !validDates.some(validDate => {
          return (
            validDate.getDate() === date.getDate() &&
            validDate.getMonth() === date.getMonth() &&
            validDate.getFullYear() === date.getFullYear()
          );
        })
    );
  };

  protected readonly contains = (d: Date, arr: DateModel[]) => {
    const find = arr.find(obj => {
      const t = moment.tz(Number(obj.date), TIMEZONE).toDate();
      return (
        t.getDate() === d.getDate() &&
        t.getMonth() === d.getMonth() &&
        t.getFullYear() === d.getFullYear()
      );
    });
    if (!find) return undefined;
    return find.times.map(a => moment.tz(Number(a), TIMEZONE).format('h:mm a'));
  };

  protected readonly onSelectedCalendarDay = (selected: Date) => {
    this.selected.next(selected);
  };

  protected readonly onSelectedCalendarMonth = (
    event: DatePickerYearChangeEvent
  ) => this.monthYearImpl(event.month, event.year);

  protected readonly onSelectedCalendarYear = (
    event: DatePickerYearChangeEvent
  ) => this.monthYearImpl(event.month, event.year);

  private monthYearImpl(month: number | undefined, year: number | undefined) {
    if (!month || !year) return;
    const d = new Date(year, month - 1, this.selected.getValue().getDate());
    this.selected.next(d);
  }
}

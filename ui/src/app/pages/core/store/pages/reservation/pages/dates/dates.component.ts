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
import { ApiResponse, ApiState } from '@root/app.model';
import {
  DateModel,
  filterValidDatesFromDatesInAMonth,
  findDatesInDateModel
} from '@shared/model/shared.model';
import { FormsModule } from '@angular/forms';
import {
  DatePicker,
  DatePickerMonthChangeEvent,
  DatePickerYearChangeEvent
} from 'primeng/datepicker';
import { ReservationModel } from '@store/pages/reservation/reservation.model';
import { Skeleton } from 'primeng/skeleton';
import { Message } from 'primeng/message';
import moment from 'moment-timezone';
import { Router } from '@angular/router';
import { RESERVATION_ROUTE } from '@store/store.routes';
import { CONFIRM_ROUTE } from '@store/pages/reservation/reservation.routes';
import { Button } from 'primeng/button';
import { PrimeTemplate } from 'primeng/api';
import { TIMEZONE } from '@root/app.util';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';
import { SummaryHolderComponent } from '@store/pages/reservation/shared/summary-holder.component';

@Component({
  selector: 'app-dates',
  imports: [
    DatePicker,
    FormsModule,
    PrimeTemplate,
    Skeleton,
    Message,
    Button,
    SummaryHolderComponent
  ],
  templateUrl: './dates.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DatesComponent {
  protected readonly reservationService = inject(ReservationService);
  protected readonly service = inject(DatesService);
  private readonly router = inject(Router);

  protected readonly state = ApiState;
  protected readonly today = new Date();
  protected date: Date[] | undefined;
  protected readonly tz = TIMEZONE;

  protected readonly selected = new BehaviorSubject<Date>(new Date());

  protected readonly skeleton = (len: number) =>
    Array.from({ length: len }, (_, index) => index);

  protected readonly validDates = signal<Date[]>([]);
  protected readonly invalidDates = signal<Date[]>([]);
  protected readonly req = toSignal(
    this.selected.asObservable().pipe(
      distinctUntilChanged(),
      debounceTime(800),
      switchMap(date =>
        this.service.dates(date).pipe(
          tap(obj => {
            if (obj.state === ApiState.LOADED && obj.data) {
              this.validDates.set(
                obj.data.map(d => moment.tz(Number(d.date), TIMEZONE).toDate())
              );
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

  /**
   * Filters out valid dates from days in the particular month
   * */
  private readonly filter = (date: Date, valid: DateModel[]) =>
    filterValidDatesFromDatesInAMonth(date, valid);

  protected readonly crossDates = (date: any, dates: Date[]) =>
    dates.some(
      d =>
        d.getDate() === date.day &&
        d.getMonth() === date.month &&
        d.getFullYear() === date.year
    );

  protected readonly contains = (d: Date, arr: DateModel[]) =>
    findDatesInDateModel(d, arr);

  protected readonly onSelectedCalendarMonth = (
    event: DatePickerMonthChangeEvent
  ) => this.monthYearImpl(event.month, event.year);

  protected readonly onSelectedCalendarYear = (
    event: DatePickerYearChangeEvent
  ) => this.monthYearImpl(event.month, event.year);

  private readonly monthYearImpl = (
    month: number | undefined,
    year: number | undefined
  ) =>
    !month || !year ? {} : this.selected.next(new Date(year, month - 1, 1));

  protected readonly updateParent = async (dt: string) => {
    this.reservationService.setDateTime(dt);
    await this.router.navigate([
      `${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${CONFIRM_ROUTE}`
    ]);
  };
}

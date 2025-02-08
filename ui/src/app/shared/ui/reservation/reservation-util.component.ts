import { ChangeDetectionStrategy, Component } from '@angular/core';
import {
  debounceTime,
  distinctUntilChanged,
  EMPTY,
  map,
  Subject,
  switchMap
} from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import {
  ConfirmModel,
  DateModel,
  StaffModel
} from '@shared/model/shared.model';
import { ServiceTypeImpl } from '@shared/data-access/service-type.service';
import { BookingService } from '@shared/data-access/booking.service';

@Component({
  selector: 'app-reservation-util',
  imports: [],
  template: ``,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ReservationUtilComponent {
  constructor(
    private readonly st: ServiceTypeImpl,
    private readonly bs: BookingService
  ) {}

  protected readonly servicesEmitter = new Subject<boolean>();
  protected readonly services = toSignal<
    ApiResponse<string[]>,
    ApiResponse<string[]>
  >(
    this.servicesEmitter.asObservable().pipe(
      switchMap(b =>
        !b
          ? EMPTY
          : this.st.services().pipe(
              map(
                o =>
                  <ApiResponse<string[]>>{
                    state: o.state,
                    message: o.message,
                    data: o.data?.map(s => s.name)
                  }
              )
            )
      )
    ),
    { initialValue: { state: ApiState.LOADED } }
  );

  protected readonly staffEmitter = new Subject<string[]>();
  protected readonly staffsByServiceTypes = toSignal(
    this.staffEmitter
      .asObservable()
      .pipe(switchMap(a => this.bs.staffsByServiceTypes(a))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<StaffModel[]> }
  );

  protected readonly validDatesEmitter = new Subject<{
    date: Date;
    services: string[];
    staff_id: string;
  }>();
  protected readonly validDates = toSignal(
    this.validDatesEmitter.asObservable().pipe(
      distinctUntilChanged(),
      debounceTime(700),
      switchMap(o => this.bs.validDates(o))
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<DateModel[]> }
  );

  protected readonly reserveBookingEmitter = new Subject<ConfirmModel>();
  protected readonly reserveBooking = toSignal(
    this.reserveBookingEmitter
      .asObservable()
      .pipe(switchMap(o => this.bs.create(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}

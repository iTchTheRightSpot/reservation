import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { BookingsService } from './bookings.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Badge } from 'primeng/badge';
import { Button } from 'primeng/button';
import { DatePicker } from 'primeng/datepicker';
import { FloatLabel } from 'primeng/floatlabel';
import { TableModule } from 'primeng/table';
import {
  BookingsModel,
  BookingsRequestPayload,
  BookingStatus,
  UpdateBookingStatusPayload
} from './bookings.model';
import { AuthService } from '@shared/data-access/auth.service';
import { EMPTY, interval, map, Subject, switchMap, takeWhile, tap } from 'rxjs';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { Avatar } from 'primeng/avatar';
import { Drawer } from 'primeng/drawer';
import { BookingDetailComponent } from './ui/detail/booking-detail.component';
import { CreateBookingComponent } from './ui/create-booking/create-booking.component';
import {
  ConfirmModel,
  DateModel,
  StaffModel
} from '@shared/model/shared.model';
import { CRMServiceTypeService } from '@crm/pages/service-type/crm-service-type.service';
import { AsyncPipe } from '@angular/common';

@Component({
  selector: 'app-bookings',
  imports: [
    Badge,
    Button,
    DatePicker,
    FloatLabel,
    TableModule,
    Select,
    FormsModule,
    Avatar,
    Drawer,
    BookingDetailComponent,
    CreateBookingComponent,
    AsyncPipe
  ],
  templateUrl: './bookings.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BookingsComponent {
  constructor() {
    interval(1000)
      .pipe(
        tap(() => {
          const u = this.authService.activeUser();
          if (u)
            this.allBookingsEmitter.next({
              date: new Date(),
              user_id: u.user_id,
              page: this.first,
              size: this.rows
            });
        }),
        takeWhile(() => !this.authService.activeUser())
      )
      .subscribe();
  }

  private readonly bookingService = inject(BookingsService);
  protected readonly authService = inject(AuthService);
  private readonly serviceTypes = inject(CRMServiceTypeService);

  protected toggleBookingDetails = false;
  protected toggleNewBookings = false;
  protected first = 0;
  protected rows = 5;
  protected readonly apiState = ApiState;
  protected readonly bookingState = BookingStatus;
  protected readonly thead = ['Name', 'From', 'To', 'Status'];

  protected selectedDate: Date | undefined;
  protected selectedStaff: CRMStaffModel | undefined;
  protected selectedBooking: BookingsModel | undefined;

  private readonly allBookingsEmitter = new Subject<BookingsRequestPayload>();
  protected readonly bookings = toSignal(
    this.allBookingsEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.bookings(o))),
    {
      initialValue: { state: ApiState.LOADED } as ApiResponse<BookingsModel[]>
    }
  );

  protected readonly onSelectedDate = (d: Date) =>
    this.allBookingsEmitter.next({
      user_id:
        this.selectedStaff?.user_id ||
        this.authService.activeUser()?.user_id ||
        '',
      date: d,
      page: (this.first = 0),
      size: this.rows
    });

  private readonly updateBookingStatusEmitter =
    new Subject<UpdateBookingStatusPayload>();
  protected readonly updateBookingStatus$ = this.updateBookingStatusEmitter
    .asObservable()
    .pipe(switchMap(o => this.bookingService.updateBookingStatus(o)));

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

  protected readonly updateBookingStatus = (o: UpdateBookingStatusPayload) => {
    this.updateBookingStatusEmitter.next(o);
    this.allBookingsEmitter.next({
      user_id:
        this.selectedStaff?.user_id ||
        this.authService.activeUser()?.user_id ||
        '',
      date: this.selectedDate || new Date(),
      page: this.first,
      size: this.rows
    });
  };

  protected readonly onSelectedStaff = () =>
    this.allBookingsEmitter.next({
      user_id: this.selectedStaff?.user_id || '',
      date: this.selectedDate || new Date(),
      page: (this.first = 0),
      size: this.rows
    });

  protected readonly servicesEmitter = new Subject<boolean>();
  protected readonly services = toSignal<
    ApiResponse<string[]>,
    ApiResponse<string[]>
  >(
    this.servicesEmitter.asObservable().pipe(
      switchMap(b =>
        !b
          ? EMPTY
          : this.serviceTypes.all().pipe(
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
  protected readonly staffs = toSignal(
    this.staffEmitter
      .asObservable()
      .pipe(switchMap(a => this.bookingService.staffsByServiceTypes(a))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<StaffModel[]> }
  );

  protected readonly validDatesEmitter = new Subject<{
    date: Date;
    services: string[];
    staff_id: string;
  }>();
  protected readonly validDates = toSignal(
    this.validDatesEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.validDates(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<DateModel[]> }
  );

  protected readonly reserveBookingEmitter = new Subject<ConfirmModel>();
  protected readonly reserveBooking = toSignal(
    this.reserveBookingEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.create(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}

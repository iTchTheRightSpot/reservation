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
import { interval, Subject, switchMap, takeWhile, tap } from 'rxjs';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { Avatar } from 'primeng/avatar';
import { Drawer } from 'primeng/drawer';
import { BookingDetailComponent } from './ui/detail/booking-detail.component';
import { CreateBookingComponent } from '@shared/ui/reservation/create-booking.component';
import { AsyncPipe } from '@angular/common';
import { ReservationUtilComponent } from '@shared/ui/reservation/reservation-util.component';
import { ServiceTypeImpl } from '@shared/data-access/service-type.service';
import { BookingService } from '@shared/data-access/booking.service';
import { CRMStaffsService } from '@crm/pages/staff/pages/all/crm-staff.service';

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
export class BookingsComponent extends ReservationUtilComponent {
  constructor(st: ServiceTypeImpl, bs: BookingService) {
    super(st, bs);

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

  private readonly bookingsService = inject(BookingsService);
  private readonly authService = inject(AuthService);
  protected readonly staffs = toSignal(inject(CRMStaffsService).staffs(), {
    initialValue: { state: ApiState.LOADED } as ApiResponse<CRMStaffModel[]>
  });

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
      .pipe(switchMap(o => this.bookingsService.bookings(o))),
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
    .pipe(switchMap(o => this.bookingsService.updateBookingStatus(o)));

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
}

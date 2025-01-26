import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { BookingsService } from './bookings.service';
import { CRMStaffsService } from '@crm/pages/staff/crm-staffs.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { CRMStaffModel } from '@crm/pages/staff/crm-staff.model';
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
import { Dialog } from 'primeng/dialog';
import { CreateBookingComponent } from './ui/create-booking/create-booking.component';
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
    Dialog,
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
  private readonly staffService = inject(CRMStaffsService);
  protected readonly authService = inject(AuthService);

  protected toggleBookingDetails = false;
  protected toggleNewBookings = false;
  protected first = 0;
  protected rows = 10;
  protected readonly apiState = ApiState;
  protected readonly bookingState = BookingStatus;
  protected readonly thead = ['Name', 'From', 'To', 'Status'];

  protected date = new Date();
  protected selectedStaff: CRMStaffModel | undefined;
  protected selectedBooking: BookingsModel | undefined;

  protected readonly staffs = toSignal(this.staffService.staffs(), {
    initialValue: { state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>
  });

  private readonly allBookingsEmitter = new Subject<BookingsRequestPayload>();
  protected readonly bookings = toSignal(
    this.allBookingsEmitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.bookings(o))),
    {
      initialValue: { state: ApiState.LOADING } as ApiResponse<BookingsModel[]>
    }
  );

  protected readonly selectedDate = (d: Date) => {
    this.date = d;
    if (!this.selectedStaff)
      this.allBookingsEmitter.next({
        user_id: this.authService.activeUser()?.user_id || '',
        date: d,
        page: (this.first = 0),
        size: this.rows
      });
    else
      this.allBookingsEmitter.next({
        user_id: this.selectedStaff.user_id,
        date: d,
        page: (this.first = 0),
        size: this.rows
      });
  };

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
    if (this.selectedStaff)
      this.allBookingsEmitter.next({
        user_id: this.selectedStaff.user_id,
        date: this.date,
        page: this.first,
        size: this.rows
      });
    else
      this.allBookingsEmitter.next({
        user_id: this.authService.activeUser()?.user_id || '',
        date: this.date,
        page: this.first,
        size: this.rows
      });
  };
}

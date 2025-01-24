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
  BookingStatus
} from './bookings.model';
import { AuthService } from '@shared/data-access/auth.service';
import { interval, Subject, switchMap, takeWhile, tap } from 'rxjs';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { Avatar } from 'primeng/avatar';
import { Drawer } from 'primeng/drawer';

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
    Drawer
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
            this.emitter.next({
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

  protected toggleDetails = false
  protected first = 0;
  protected rows = 10;
  protected date = new Date();
  protected readonly apiState = ApiState;
  protected readonly bookingState = BookingStatus;
  protected readonly thead = ['Name', 'From', 'To', 'Status'];

  protected selectedStaff: CRMStaffModel | undefined;

  protected readonly staffs = toSignal(this.staffService.staffs(), {
    initialValue: { state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>
  });

  private readonly emitter = new Subject<BookingsRequestPayload>();
  protected readonly bookings = toSignal(
    this.emitter
      .asObservable()
      .pipe(switchMap(o => this.bookingService.bookings(o))),
    {
      initialValue: { state: ApiState.LOADING } as ApiResponse<BookingsModel[]>
    }
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
}

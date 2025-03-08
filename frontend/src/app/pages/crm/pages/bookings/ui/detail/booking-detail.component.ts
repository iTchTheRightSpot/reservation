import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import {
  BookingsModel,
  BookingStatus,
  UpdateBookingStatusPayload
} from '@crm/pages/bookings/bookings.model';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Textarea } from 'primeng/textarea';
import { Select } from 'primeng/select';
import { Divider } from 'primeng/divider';
import { InputNumber } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { Tag } from 'primeng/tag';
import { Button } from 'primeng/button';
import { ConfirmPopup } from 'primeng/confirmpopup';
import { ConfirmationService } from 'primeng/api';
import { ApiState } from '@root/app.model';
import { Avatar } from 'primeng/avatar';

@Component({
  selector: 'app-booking-detail',
  imports: [
    IconField,
    InputIcon,
    InputText,
    Textarea,
    Select,
    Divider,
    InputNumber,
    FormsModule,
    Tag,
    Button,
    ConfirmPopup,
    Avatar
  ],
  providers: [ConfirmationService],
  templateUrl: './booking-detail.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BookingDetailComponent {
  private readonly confirmService = inject(ConfirmationService);

  booking = input.required<BookingsModel>();
  updateStatus = input.required<ApiState | undefined>();

  readonly emitter = output<UpdateBookingStatusPayload>();

  protected readonly state = BookingStatus;
  protected selectedBookingStatus: BookingStatus | undefined;

  protected readonly statuses = (s: BookingStatus) =>
    [BookingStatus.CONFIRMED, BookingStatus.CANCELLED]
      .filter(a => a !== s)
      .map(a => a);

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

  protected readonly confirm = (event: Event, reservationId: number) =>
    this.confirmService.confirm({
      target: event.target as EventTarget,
      message: 'Are you sure you want to proceed?',
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptButtonProps: {
        label: 'Save'
      },
      accept: () => {
        const s = this.selectedBookingStatus;
        if (s) this.emitter.emit({ reservation_id: reservationId, status: s });
      }
    });
}

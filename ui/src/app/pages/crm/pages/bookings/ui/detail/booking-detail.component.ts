import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import {
  BookingsModel,
  BookingStatus
} from '@crm/pages/bookings/bookings.model';
import { FloatLabel } from 'primeng/floatlabel';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Textarea } from 'primeng/textarea';
import { Badge } from 'primeng/badge';
import { Select } from 'primeng/select';
import { Divider } from 'primeng/divider';

@Component({
  selector: 'app-booking-detail',
  imports: [
    FloatLabel,
    IconField,
    InputIcon,
    InputText,
    Textarea,
    Badge,
    Select,
    Divider
  ],
  templateUrl: './booking-detail.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BookingDetailComponent {
  booking = input.required<BookingsModel>();
  protected readonly state = BookingStatus;

  protected readonly services = (arr: string[]) =>
    arr.map(s => ({ name: s.trim() }));
}

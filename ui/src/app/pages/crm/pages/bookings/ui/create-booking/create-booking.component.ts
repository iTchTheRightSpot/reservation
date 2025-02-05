import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Button } from 'primeng/button';
import { FloatLabel } from 'primeng/floatlabel';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Message } from 'primeng/message';
import {
  FormBuilder,
  FormControl,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { Textarea } from 'primeng/textarea';
import { ApiResponse, ApiState } from '@root/app.model';
import { Select } from 'primeng/select';
import { MultiSelect } from 'primeng/multiselect';
import { Avatar } from 'primeng/avatar';
import { Divider } from 'primeng/divider';
import { DatePicker } from 'primeng/datepicker';
import { ConfirmModel, DateModel } from '@shared/data-access/shared.model';
import { TIMEZONE } from '@root/app.util';

@Component({
  selector: 'app-create-booking',
  imports: [
    Button,
    FloatLabel,
    IconField,
    InputIcon,
    InputText,
    Message,
    ReactiveFormsModule,
    Textarea,
    Select,
    MultiSelect,
    Avatar,
    Divider,
    DatePicker
  ],
  templateUrl: './create-booking.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CreateBookingComponent {
  services = input.required<ApiResponse<string[]>>();
  staffs = input.required<ApiResponse<CRMStaffModel[]>>();
  datesTimes = input.required<ApiResponse<DateModel[]>>();
  submissionRequestState = input.required<ApiResponse<any>>();

  readonly staffEmitter = output<string[]>();
  readonly datesEmitter = output<{ services: string[]; staff_id: string }>();
  readonly submitEmitter = output<ConfirmModel>();

  protected readonly state = ApiState;
  protected readonly today = new Date();

  protected readonly form = inject(FormBuilder).group({
    phone: new FormControl<string>('', [
      Validators.required,
      Validators.minLength(9),
      Validators.maxLength(20)
    ]),
    name: new FormControl<string>('', [
      Validators.required,
      Validators.minLength(1),
      Validators.maxLength(100)
    ]),
    email: new FormControl<string>('', [
      Validators.required,
      Validators.email,
      Validators.maxLength(320)
    ]),
    description: new FormControl<string>('', [Validators.maxLength(255)]),
    services: new FormControl<string[] | null>(null, [
      Validators.required,
      Validators.min(1)
    ]),
    staff_id: new FormControl<string | null>({ value: null, disabled: true }, [
      Validators.required
    ]),
    date: new FormControl<Date | null>({ value: null, disabled: true }, [
      Validators.required
    ]),
    date_time: new FormControl<string | null>({ value: null, disabled: true }, [
      Validators.required
    ])
  });

  protected readonly submit = () =>
    this.form.invalid
      ? {}
      : this.submitEmitter.emit({
          staff_id: this.form.controls['staff_id'].value || '',
          name: this.form.controls['name'].value || '',
          email: this.form.controls['email'].value || '',
          description: this.form.controls['description'].value || '',
          phone: this.form.controls['phone'].value || '',
          services: this.form.controls['services'].value || [],
          timezone: TIMEZONE,
          time: this.form.controls['date_time'].value || ''
        });
}

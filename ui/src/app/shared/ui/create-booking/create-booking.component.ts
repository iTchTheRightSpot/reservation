import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  OnChanges,
  output,
  signal,
  SimpleChanges
} from '@angular/core';
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
import {
  DatePicker,
  DatePickerMonthChangeEvent,
  DatePickerYearChangeEvent
} from 'primeng/datepicker';
import {
  ConfirmModel,
  DateModel,
  filterValidDatesFromDatesInAMonth,
  findDatesInDateModel,
  StaffModel
} from '@shared/model/shared.model';
import { TIMEZONE } from '@root/app.util';
import moment from 'moment-timezone';
import { PrimeTemplate } from 'primeng/api';

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
    DatePicker,
    PrimeTemplate
  ],
  templateUrl: './create-booking.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CreateBookingComponent implements OnChanges {
  services = input.required<ApiResponse<string[]>>();
  staffs = input.required<ApiResponse<StaffModel[]>>();
  datesTimes = input.required<ApiResponse<DateModel[]>>();
  submissionRequestState = input.required<ApiState>();

  readonly staffEmitter = output<string[]>();
  readonly datesEmitter = output<{
    services: string[];
    staff_id: string;
    date: Date;
  }>();
  readonly submitEmitter = output<ConfirmModel>();

  protected readonly tz = TIMEZONE;
  protected readonly state = ApiState;
  protected readonly today = new Date();
  protected readonly invalidDates = signal<Date[]>([]);
  protected readonly validDates = signal<Date[]>([]);

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
    staff_id: new FormControl<StaffModel | null>(
      { value: null, disabled: true },
      [Validators.required]
    ),
    date: new FormControl<Date | null>({ value: new Date(), disabled: true }, [
      Validators.required
    ]),
    date_time: new FormControl<{ original: string; format: string } | null>(
      { value: null, disabled: true },
      [Validators.required]
    )
  });

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['datesTimes']) {
      const dts = changes['datesTimes'].currentValue as ApiResponse<
        DateModel[]
      >;
      if (dts.state === ApiState.LOADED && dts.data) {
        this.validDates.set(
          dts.data.map(d => moment.tz(Number(d.date), TIMEZONE).toDate())
        );
        this.invalidDates.set(
          this.filter(this.form.controls['date'].value!!, dts.data)
        );
      }
    }
  }

  protected readonly crossDates = (date: any, dates: Date[]) =>
    dates.some(
      d =>
        d.getDate() === date.day &&
        d.getMonth() === date.month &&
        d.getFullYear() === date.year
    );

  protected readonly contains = (d: Date, arr: DateModel[]) =>
    findDatesInDateModel(d, arr);

  private readonly filter = (date: Date, valid: DateModel[]) =>
    filterValidDatesFromDatesInAMonth(date, valid);

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
    !month || !year
      ? {}
      : this.datesEmitter.emit({
          date: new Date(year, month - 1, 1),
          staff_id: this.form.controls['staff_id'].value?.staff_id || '',
          services: this.form.controls['services'].value || []
        });

  protected readonly submit = () =>
    this.form.invalid
      ? {}
      : this.submitEmitter.emit({
          staff_id: this.form.controls['staff_id'].value?.staff_id || '',
          name: this.form.controls['name'].value || '',
          email: this.form.controls['email'].value || '',
          description: this.form.controls['description'].value || '',
          phone: this.form.controls['phone'].value || '',
          services: this.form.controls['services'].value || [],
          timezone: TIMEZONE,
          time: this.form.controls['date_time'].value?.original || ''
        });
}

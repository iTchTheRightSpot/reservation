import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { DatePicker } from 'primeng/datepicker';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { Select } from 'primeng/select';
import { InputText } from 'primeng/inputtext';
import { Button } from 'primeng/button';
import { Message } from 'primeng/message';
import { ApiState } from '@root/app.model';
import { CreateUpdateScheduleModel } from './shared-create-update-schedule.model';

@Component({
  selector: 'app-shared-create-update-schedule',
  imports: [
    DatePicker,
    FormsModule,
    Select,
    InputText,
    ReactiveFormsModule,
    Button,
    Message
  ],
  template: `
    <form [formGroup]="form()" class="w-full flex flex-col gap-4">
      <!-- calendar -->
      <div class="flex-auto">
        <label for="calendar-12h" class="font-bold block mb-2">
          Calendar
        </label>
        <p-date-picker
          styleClass="w-full"
          inputId="calendar-12h"
          formControlName="date_time"
          [showIcon]="true"
          [showTime]="true"
          [hourFormat]="'12'"
        />
        @if (
          form().controls['date_time'].invalid &&
          !form().controls['date_time'].untouched
        ) {
          <p-message severity="error" variant="simple" size="small">
            Date & Time are required
          </p-message>
        }
      </div>

      <!-- visible & reoccurring -->
      <div class="w-full flex gap-2 flex-col md:flex-row items-end">
        <div class="w-full flex-auto">
          <label class="font-semibold block mb-2" for="visible">Visible</label>
          <p-select
            formControlName="is_visible"
            id="visible"
            placeholder="allow customers to book time"
            [options]="[true, false]"
            class="w-full"
          />
          @if (
            form().controls['is_visible'].invalid &&
            !form().controls['is_visible'].untouched
          ) {
            <p-message severity="error" variant="simple" size="small">
              Visibility is required
            </p-message>
          }
        </div>

        <div class="w-full flex-auto">
          <label class="font-semibold block mb-2" for="reoccurring"
            >Reoccurring</label
          >
          <p-select
            formControlName="is_reoccurring"
            id="reoccurring"
            placeholder="automatically schedule every week"
            [options]="[true, false]"
            class="w-full"
          />
          @if (
            form().controls['is_reoccurring'].invalid &&
            !form().controls['is_reoccurring'].untouched
          ) {
            <p-message severity="error" variant="simple" size="small">
              Reoccurring is required
            </p-message>
          }
        </div>
      </div>

      <!-- duration -->
      <div class="flex-1">
        <label for="duration" class="font-bold block mb-2">Duration</label>
        <input
          formControlName="duration"
          id="duration"
          class="w-full"
          type="text"
          placeholder="seconds"
          pInputText
        />
        @if (
          form().controls['duration'].invalid &&
          !form().controls['duration'].untouched
        ) {
          <p-message severity="error" variant="simple" size="small">
            Duration is required
          </p-message>
        }
      </div>

      <p-button
        (onClick)="submit(staffId())"
        [className]="'text-right'"
        [disabled]="form().invalid || loadingState() === state.LOADING"
      >
        @if (loadingState() === state.LOADING) {
          <i class="pi pi-spin pi-spinner" style="font-size: 1rem"></i>
        } @else {
          Create
        }
      </p-button>
    </form>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedCreateUpdateScheduleComponent {
  form = input.required<FormGroup>();
  staffId = input.required<string>();
  loadingState = input.required<ApiState>();

  readonly emitter = output<CreateUpdateScheduleModel>();

  protected readonly state = ApiState;

  protected readonly submit = (staffId: string) => {
    if (this.form().invalid) return;
    this.emitter.emit({
      staff_id: staffId,
      date_time: this.form().controls['date_time'].value,
      is_visible: this.form().controls['is_visible'].value,
      is_reoccurring: this.form().controls['is_reoccurring'].value,
      duration: this.form().controls['duration'].value
    } as CreateUpdateScheduleModel);
  };
}

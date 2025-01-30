import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  OnChanges,
  output,
  SimpleChanges
} from '@angular/core';
import { ApiState } from '@root/app.model';
import {
  FormBuilder,
  FormControl,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import {
  Schedule,
  UpdateScheduleModel
} from '@crm/pages/account/pages/schedule/schedule.model';
import { Divider } from 'primeng/divider';
import { InputText } from 'primeng/inputtext';
import { Message } from 'primeng/message';
import { Select } from 'primeng/select';
import { Button } from 'primeng/button';

@Component({
  selector: 'app-edit-schedule',
  imports: [Divider, InputText, Message, ReactiveFormsModule, Select, Button],
  template: `
    <form [formGroup]="form" class="h-[30rem] w-full flex flex-col gap-4">
      <!-- from & to -->
      <div
        class="w-full flex flex-col md:flex-row items-center justify-center gap-2"
      >
        <div class="w-full">
          <div class="flex-auto">
            <label class="block" for="from">From</label>
            <input
              id="from"
              [disabled]="true"
              [value]="fm(schedule().start)"
              type="text"
              class="w-full"
              pInputText
              placeholder="name*"
            />
          </div>
        </div>

        <div class="w-full">
          <div class="flex-auto">
            <label class="block" for="expire">To</label>
            <input
              [disabled]="true"
              [value]="fm(schedule().end)"
              type="text"
              id="expire"
              class="w-full"
              pInputText
              placeholder="expire at"
            />
          </div>
        </div>
      </div>

      <p-divider />

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
            form.controls['is_visible'].invalid &&
            !form.controls['is_visible'].untouched
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
            form.controls['is_reoccurring'].invalid &&
            !form.controls['is_reoccurring'].untouched
          ) {
            <p-message severity="error" variant="simple" size="small">
              Reoccurring is required
            </p-message>
          }
        </div>
      </div>

      <p-button
        (onClick)="emit(staffId(), schedule().schedule_id)"
        [className]="'text-right'"
        [disabled]="form.invalid || loadingState() === state.LOADING"
      >
        @if (loadingState() === state.LOADING) {
          <i class="pi pi-spin pi-spinner" style="font-size: 1rem"></i>
        } @else {
          Edit
        }
      </p-button>
    </form>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class EditScheduleComponent implements OnChanges {
  schedule = input.required<Schedule>();
  staffId = input.required<string>();
  loadingState = input.required<ApiState>();

  readonly updateScheduleEmitter = output<UpdateScheduleModel>();

  protected readonly state = ApiState;
  protected readonly form = inject(FormBuilder).group({
    is_visible: new FormControl<boolean>(false, [Validators.required]),
    is_reoccurring: new FormControl<boolean>(false, [Validators.required])
  });

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['schedule']) {
      this.form.controls['is_visible'].setValue(this.schedule().is_visible);
      this.form.controls['is_reoccurring'].setValue(
        this.schedule().is_reoccurring
      );
    }
  }

  protected readonly emit = (staffId: string, scheduleId: number) =>
    this.updateScheduleEmitter.emit(<UpdateScheduleModel>{
      schedule_id: scheduleId,
      staff_id: staffId,
      is_reoccurring: this.form.controls['is_reoccurring'].value,
      is_visible: this.form.controls['is_visible'].value
    });

  protected readonly fm = (d: string) =>
    new Date(Number(d)).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour12: true,
      hour: 'numeric',
      minute: 'numeric'
    });
}

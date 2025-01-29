import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { SharedCreateUpdateScheduleComponent } from './shared/shared-create-update-schedule.component';
import { ApiState } from '@root/app.model';
import { CreateUpdateScheduleModel } from './shared/shared-create-update-schedule.model';
import { FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-create-schedule',
  imports: [SharedCreateUpdateScheduleComponent],
  template: `
    <app-shared-create-update-schedule
      [form]="form"
      [staffId]="staffId()"
      [loadingState]="loadingState()"
      (emitter)="createScheduleEmitter.emit($event)"
    />
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CreateScheduleComponent {
  protected readonly form = inject(FormBuilder).group({
    date_time: new FormControl<Date | null>(null, [Validators.required]),
    is_visible: new FormControl<boolean>(false, [Validators.required]),
    is_reoccurring: new FormControl<boolean>(false, [Validators.required]),
    duration: new FormControl<number | null>(null, [Validators.required])
  });

  staffId = input.required<string>();
  loadingState = input.required<ApiState>();

  readonly createScheduleEmitter = output<CreateUpdateScheduleModel>();
}

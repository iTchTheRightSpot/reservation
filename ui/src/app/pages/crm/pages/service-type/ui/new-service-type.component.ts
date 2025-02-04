import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  model,
  output
} from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { CRMServiceTypeModel } from '@crm/pages/service-type/crm-service-type.model';
import { ApiState } from '@root/app.model';
import { ServiceTypeFormComponent } from './shared/service-type-form.component';
import { ServiceTypeSelectVisibilityModel } from './shared.model';

@Component({
  selector: 'app-new-service-type',
  imports: [ServiceTypeFormComponent],
  template: `
    <app-service-type-form
      [form]="form"
      [isLoading]="apiState() === state.LOADING"
      [visibilityOptions]="arr"
      (submit)="sub()"
      (cancel)="visible.set(false)"
    />
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NewServiceTypeComponent {
  apiState = input.required<ApiState>();
  visible = model.required<boolean>();
  readonly emitter = output<CRMServiceTypeModel>();

  protected readonly state = ApiState;
  protected readonly arr: ServiceTypeSelectVisibilityModel[] = [
    { name: 'True', value: true },
    { name: 'False', value: false }
  ];

  protected form = inject(FormBuilder).group({
    name: new FormControl<string>('', [
      Validators.required,
      Validators.minLength(1),
      Validators.maxLength(50)
    ]),
    price: new FormControl<number | null>(null, [Validators.required]),
    is_visible: new FormControl<ServiceTypeSelectVisibilityModel | null>(null, [
      Validators.required
    ]),
    duration: new FormControl<number | null>(null, [Validators.required]),
    clean_up_time: new FormControl<number | null>(null, [Validators.required])
  });

  protected readonly sub = () => {
    if (this.form.invalid) return;
    this.emitter.emit({
      name: this.form.value.name!!,
      price: Number(this.form.value.price!!),
      is_visible: this.form.value.is_visible!!.value,
      duration: Number(this.form.value.duration!!),
      clean_up_time: Number(this.form.value.clean_up_time!!)
    } as CRMServiceTypeModel);
  };
}

import {
  ChangeDetectionStrategy,
  Component,
  input,
  OnChanges,
  SimpleChanges
} from '@angular/core';
import { CRMServiceTypeModel } from '@crm/pages/service-type/crm-service-type.model';
import { NewServiceTypeComponent } from './new-service-type.component';
import { ServiceTypeFormComponent } from './shared/service-type-form.component';

@Component({
  selector: 'app-edit-service-type',
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
export class EditServiceTypeComponent
  extends NewServiceTypeComponent
  implements OnChanges
{
  service = input.required<CRMServiceTypeModel | undefined>();

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['service']) {
      const curr = changes['service'].currentValue as CRMServiceTypeModel;
      if (curr) {
        this.form.controls['name'].setValue(curr.name);
        this.form.controls['price'].setValue(curr.price);
        this.form.controls['is_visible'].setValue({
          name: `${curr.is_visible}`,
          value: curr.is_visible
        });
        this.form.controls['duration'].setValue(curr.duration);
        this.form.controls['clean_up_time'].setValue(curr.clean_up_time);
      }
    }
  }
}

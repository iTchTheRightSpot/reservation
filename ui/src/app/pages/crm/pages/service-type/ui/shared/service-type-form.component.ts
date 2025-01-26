import {
  ChangeDetectionStrategy,
  Component,
  input,
  output
} from '@angular/core';
import { FormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { ServiceTypeSelectVisibilityModel } from '@crm/pages/service-type/ui/shared.model';
import { KeyFilter } from 'primeng/keyfilter';
import { InputNumber } from 'primeng/inputnumber';

@Component({
  selector: 'app-service-type-form',
  imports: [
    Button,
    FormsModule,
    InputText,
    ReactiveFormsModule,
    Select,
    KeyFilter,
    InputNumber
  ],
  template: `
    <form [formGroup]="form()">
      <div class="flex items-center gap-4 mb-4">
        <label for="name" class="font-semibold w-24">Name</label>
        <input
          pInputText
          id="name"
          class="flex-auto"
          formControlName="name"
          autocomplete="off"
        />
      </div>
      <div class="flex items-center gap-4 mb-8">
        <label for="price" class="font-semibold w-24">Price (₦)</label>
        <p-input-number
          id="price"
          class="flex-1"
          formControlName="price"
          inputId="minmaxfraction"
          mode="decimal"
          [minFractionDigits]="2"
          [maxFractionDigits]="2"
        />
      </div>
      <div class="flex items-center gap-4 mb-8">
        <label for="duration" class="font-semibold w-24"
          >Duration(in seconds)</label
        >
        <input
          pInputText
          id="duration"
          pKeyFilter="int"
          class="flex-auto"
          formControlName="duration"
          autocomplete="off"
        />
      </div>
      <div class="flex items-center gap-4 mb-8">
        <label for="clean_up" class="font-semibold w-24"
          >Clean up(in seconds)</label
        >
        <input
          pInputText
          id="clean_up"
          pKeyFilter="int"
          formControlName="clean_up_time"
          class="flex-auto"
          autocomplete="off"
        />
      </div>
      <div class="flex items-center gap-4 mb-8">
        <label for="is_visible" class="font-semibold w-24">Visible</label>
        <p-select
          inputId="is_visible"
          [options]="visibilityOptions()"
          formControlName="is_visible"
          optionLabel="name"
          class="flex-1"
          [placeholder]="
            form().controls['is_visible'].value?.name ||
            'Display in store front'
          "
        />
      </div>

      <div class="flex justify-end gap-2">
        <p-button label="Cancel" severity="secondary" (click)="cancel.emit()" />
        <p-button
          label="Save"
          (click)="submit.emit()"
          [disabled]="form().invalid || isLoading()"
        />
      </div>
    </form>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ServiceTypeFormComponent {
  form = input.required<FormGroup>();
  isLoading = input.required<boolean>();
  visibilityOptions = input.required<ServiceTypeSelectVisibilityModel[]>();

  readonly cancel = output<void>();
  readonly submit = output<void>();
}

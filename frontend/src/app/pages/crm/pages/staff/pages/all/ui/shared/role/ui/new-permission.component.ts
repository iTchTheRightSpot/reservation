import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { ApiState } from '@root/app.model';
import { Permission, Role } from '@crm/pages/staff/pages/all/crm-staff.model';
import {
  FormBuilder,
  FormControl,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { Button } from 'primeng/button';
import { MultiSelect } from 'primeng/multiselect';
import { RoleAndPermissionPayload } from '@crm/pages/account/account.model';

@Component({
  selector: 'app-new-permission',
  imports: [Button, MultiSelect, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" class="w-full flex flex-col gap-3">
      <p-multi-select
        appendTo="body"
        formControlName="permissions"
        [options]="permissions()"
        placeholder="select permissions"
        styleClass="w-full"
      />

      <p-button
        type="button"
        (onClick)="submit(staffId(), role())"
        [disabled]="form.invalid || addPermissionState() === state.LOADING"
      >
        @if (addPermissionState() === state.LOADING) {
          <i class="pi pi-spin pi-spinner" style="font-size: 1rem"></i>
        } @else {
          Submit
        }
      </p-button>
    </form>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NewPermissionComponent {
  addPermissionState = input.required<ApiState>();
  staffId = input.required<string>();
  role = input.required<Role>();
  permissions = input.required<Permission[]>();

  readonly emitter = output<RoleAndPermissionPayload>();

  protected readonly state = ApiState;

  protected readonly form = inject(FormBuilder).group({
    permissions: new FormControl<Permission[] | null>(null, [
      Validators.required
    ])
  });

  protected readonly submit = (staffId: string, role: Role) =>
    this.form.invalid
      ? {}
      : this.emitter.emit({
          user_id: staffId,
          role_permission: [
            {
              role: role,
              permissions: this.form.controls['permissions'].value!!
            }
          ]
        });
}

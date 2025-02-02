import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { Permission, Role } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Select } from 'primeng/select';
import {
  FormBuilder,
  FormControl,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { MultiSelect } from 'primeng/multiselect';
import { Button } from 'primeng/button';
import { RoleAndPermissionPayload } from '@crm/pages/account/account.model';
import { ApiState } from '@root/app.model';

@Component({
  selector: 'app-new-role-permission',
  imports: [Select, MultiSelect, ReactiveFormsModule, Button],
  template: `
    <form [formGroup]="form" class="w-full flex flex-col gap-3">
      <p-select
        appendTo="body"
        [options]="roles()"
        (onChange)="
          form.controls['role'].value === null
            ? form.controls['permissions'].disable()
            : form.controls['permissions'].enable()
        "
        formControlName="role"
        placeholder="select a role"
        class="w-full"
      />
      <p-multi-select
        appendTo="body"
        [options]="permissions"
        formControlName="permissions"
        placeholder="select permissions"
        styleClass="w-full"
      />

      <p-button
        type="button"
        [disabled]="
          form.invalid || createRolesAndPermissionState() === state.LOADING
        "
        (onClick)="submit(staffId())"
      >
        @if (createRolesAndPermissionState() === state.LOADING) {
          <i class="pi pi-spin pi-spinner" style="font-size: 1rem"></i>
        } @else {
          Submit
        }
      </p-button>
    </form>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NewRolePermissionComponent {
  createRolesAndPermissionState = input.required<ApiState>();
  staffId = input.required<string>();
  roles = input.required<Role[]>();

  readonly emitter = output<RoleAndPermissionPayload>();

  protected readonly form = inject(FormBuilder).group({
    role: new FormControl<Role | null>(null, [Validators.required]),
    permissions: new FormControl<Permission[] | null>(null, [
      Validators.required
    ])
  });

  protected readonly state = ApiState;
  protected readonly permissions = [
    Permission.READ,
    Permission.WRITE,
    Permission.DELETE
  ];

  protected readonly submit = (staffId: string) =>
    this.form.invalid
      ? {}
      : this.emitter.emit({
          user_id: staffId,
          role_permission: [
            {
              role: this.form.controls['role'].value!!,
              permissions: this.form.controls['permissions'].value!!
            }
          ]
        });
}

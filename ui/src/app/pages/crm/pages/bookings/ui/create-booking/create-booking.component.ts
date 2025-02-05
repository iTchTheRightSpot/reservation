import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input
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
import { Tab, TabList, TabPanel, TabPanels, Tabs } from 'primeng/tabs';
import { Divider } from 'primeng/divider';

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
    Tab,
    TabList,
    TabPanel,
    TabPanels,
    Tabs,
    Divider
  ],
  templateUrl: './create-booking.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CreateBookingComponent {
  staffs = input.required<CRMStaffModel[]>();
  request = input.required<ApiResponse<any>>();

  services = input.required<string[]>();
  protected readonly state = ApiState;

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
    staff_id: new FormControl<string | null>(null, [Validators.required]),
    services: new FormControl<string[] | null>(null, [
      Validators.required,
      Validators.min(1)
    ])
  });
}

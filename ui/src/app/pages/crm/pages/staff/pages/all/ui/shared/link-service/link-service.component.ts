import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  output
} from '@angular/core';
import { Button } from 'primeng/button';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { ServiceTypeToStaffPayload } from './link-service.model';
import { ApiResponse, ApiState } from '@root/app.model';
import { Divider } from 'primeng/divider';
import { Skeleton } from 'primeng/skeleton';
import { ConfirmPopup } from 'primeng/confirmpopup';
import { ConfirmationService } from 'primeng/api';

@Component({
  selector: 'app-link-service',
  imports: [Button, Select, FormsModule, Divider, Skeleton, ConfirmPopup],
  providers: [ConfirmationService],
  templateUrl: './link-service.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LinkServiceComponent {
  private readonly confirmService = inject(ConfirmationService);

  staffId = input.required<string>();
  allServices = input.required<ApiResponse<string[]>>();
  staffServices = input.required<ApiResponse<string[]>>();
  linkServiceLoadingState = input.required<ApiState>();
  deLinkServiceFromStaffState = input.required<ApiState>();

  readonly deLinkServiceFromStaffEmitter = output<ServiceTypeToStaffPayload>();
  readonly linkServiceToStaffEmitter = output<ServiceTypeToStaffPayload>();

  protected selectedService: string | undefined;
  protected readonly state = ApiState;

  protected readonly filter = (
    allServices: string[],
    staffServices: string[]
  ) => allServices.filter(a => !staffServices.includes(a));

  protected readonly confirm = (
    event: Event,
    service: string,
    staffId: string
  ) =>
    this.confirmService.confirm({
      target: event.target as EventTarget,
      message: `confirm de-linking of ${service}`,
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true
      },
      acceptButtonProps: { label: 'de-link' },
      accept: () =>
        this.deLinkServiceFromStaffEmitter.emit({
          staff_id: staffId,
          service: service
        })
    });
}

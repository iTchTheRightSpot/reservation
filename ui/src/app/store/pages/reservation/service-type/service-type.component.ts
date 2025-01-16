import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ServiceTypeService } from './service-type.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.util';
import { Skeleton } from 'primeng/skeleton';
import { Message } from 'primeng/message';
import { ServiceTypeModel } from './service-type.model';

@Component({
  selector: 'app-service-type',
  imports: [Skeleton, Message],
  templateUrl: './service-type.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ServiceTypeComponent {
  private readonly service = inject(ServiceTypeService);

  protected readonly state = ApiState;
  protected readonly services = toSignal(this.service.services(), {
    initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
      ServiceTypeModel[]
    >
  });
}

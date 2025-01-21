import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Button } from 'primeng/button';
import { TableModule, TablePageEvent } from 'primeng/table';
import { CRM_ServiceTypeModel } from './crm-service-type.model';
import { Badge } from 'primeng/badge';
import { CRMServiceTypeService } from './crm-service-type.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';

@Component({
  selector: 'app-crm-service-type',
  imports: [Button, TableModule, Badge],
  templateUrl: './crm-service-type.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CRMServiceTypeComponent {
  private readonly service = inject(CRMServiceTypeService);

  protected first = 0;
  protected rows = 10;

  protected readonly models = toSignal(this.service.all(), {
    initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
      CRM_ServiceTypeModel[]
    >
  });

  protected readonly pageChange = (event: TablePageEvent) => {
    this.first = event.first;
    this.rows = event.rows;
  };
}

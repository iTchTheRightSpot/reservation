import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal
} from '@angular/core';
import { Button } from 'primeng/button';
import { TableModule, TablePageEvent } from 'primeng/table';
import { CRM_ServiceTypeModel } from './crm-service-type.model';
import { Badge } from 'primeng/badge';
import { CRMServiceTypeService } from './crm-service-type.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { Dialog } from 'primeng/dialog';
import { NewServiceTypeComponent } from './ui/new-service-type.component';
import { EditServiceTypeComponent } from './ui/edit-service-type.component';
import { Skeleton } from 'primeng/skeleton';
import { Subject, switchMap } from 'rxjs';

@Component({
  selector: 'app-crm-service-type',
  imports: [
    Button,
    TableModule,
    Badge,
    Dialog,
    NewServiceTypeComponent,
    EditServiceTypeComponent,
    Skeleton
  ],
  templateUrl: './crm-service-type.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CRMServiceTypeComponent {
  protected readonly service = inject(CRMServiceTypeService);

  protected first = 0;
  protected rows = 10;
  protected readonly state = ApiState;
  protected readonly thead = [
    'Name',
    'Price',
    'Visible',
    'Duration',
    'Clean up'
  ];
  protected toggleNewServiceTypeView = false;
  protected toggleEditServiceTypeView = false;

  protected readonly clickedServiceType = signal<
    CRM_ServiceTypeModel | undefined
  >(undefined);

  protected readonly models = toSignal(this.service.all(), {
    initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
      CRM_ServiceTypeModel[]
    >
  });

  protected readonly create = new Subject<CRM_ServiceTypeModel>();
  protected readonly createApiState = toSignal(
    this.create.asObservable().pipe(switchMap(o => this.service.create(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly update = new Subject<CRM_ServiceTypeModel>();
  protected readonly updateApiState = toSignal(
    this.update.asObservable().pipe(switchMap(o => this.service.update(o))),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly tabData = (m: ApiResponse<CRM_ServiceTypeModel[]>) => {
    if (m.state === ApiState.LOADING)
      return Array.from({ length: 10 }).map(() => ({}));
    else if (m.state === ApiState.ERROR)
      return Array.from({ length: 1 }).map(() => ({}));
    return m.data!!;
  };

  protected readonly pageChange = (event: TablePageEvent) => {
    this.first = event.first;
    this.rows = event.rows;
  };
}

import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ServiceTypeService } from './service-type.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { Skeleton } from 'primeng/skeleton';
import { Message } from 'primeng/message';
import { ServiceTypeModel } from './service-type.model';
import { CardModule } from 'primeng/card';
import { FormsModule } from '@angular/forms';
import { RadioButton } from 'primeng/radiobutton';
import { Button } from 'primeng/button';
import { NgClass } from '@angular/common';
import { Router } from '@angular/router';
import { RESERVATION_ROUTE } from '@store/store.routes';
import { STORE_STAFF_ROUTE } from '@store/pages/reservation/reservation.routes';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { FORMAT_SECONDS } from '@root/app.util';

@Component({
  selector: 'app-service-type',
  imports: [
    Skeleton,
    Message,
    CardModule,
    FormsModule,
    RadioButton,
    Button,
    NgClass
  ],
  templateUrl: './service-type.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ServiceTypeComponent {
  constructor() {
    const s = this.service1.reservationState().services;
    if (s) this.selectedServices.push(...s);
  }

  private readonly service1 = inject(ReservationService);
  private readonly service2 = inject(ServiceTypeService);
  private readonly router = inject(Router);

  protected readonly state = ApiState;
  protected readonly services = toSignal(this.service2.services(), {
    initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
      ServiceTypeModel[]
    >
  });

  protected readonly formatSeconds = (seconds: number) =>
    FORMAT_SECONDS(seconds);

  protected readonly selectedServices: ServiceTypeModel[] = [];

  protected readonly contains = (o: ServiceTypeModel) =>
    this.selectedServices.some(s => s.name === o.name);

  protected readonly selected = (o: ServiceTypeModel) => {
    if (this.selectedServices.some(s => s.name === o.name)) {
      for (let i = 0; i < this.selectedServices.length; i += 1)
        if (this.selectedServices[i].name === o.name)
          this.selectedServices.splice(i, 1);
    } else this.selectedServices.push(o);
  };

  protected readonly submit = async () => {
    this.service1.setServiceTypes(this.selectedServices);
    await this.router.navigate([`/${RESERVATION_ROUTE}/${STORE_STAFF_ROUTE}`]);
  };
}

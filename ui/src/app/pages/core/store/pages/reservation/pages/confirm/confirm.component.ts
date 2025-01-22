import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { SummaryComponent } from '@store/pages/reservation/shared/summary.component';
import { ReservationModel } from '@store/pages/reservation/reservation.model';

@Component({
  selector: 'app-confirm',
  imports: [
    SummaryComponent
  ],
  templateUrl: './confirm.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConfirmComponent {
  protected readonly service1 = inject(ReservationService);

  protected readonly check = (o: ReservationModel) =>
    !o.services || o.services.length < 1 || !o.staff
      ? undefined
      : { services: o.services, staff: o.staff };
}

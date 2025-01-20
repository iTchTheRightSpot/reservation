import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';

@Component({
  selector: 'app-confirm',
  imports: [],
  templateUrl: './confirm.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConfirmComponent {
  protected readonly service1 = inject(ReservationService);
}

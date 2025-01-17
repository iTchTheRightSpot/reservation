import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal
} from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { Router } from '@angular/router';
import { ApiResponse, ApiState } from '@root/app.util';
import { toSignal } from '@angular/core/rxjs-interop';
import { StaffService } from './staff.service';
import { StaffModel } from './staff.model';
import { Message } from 'primeng/message';
import { Skeleton } from 'primeng/skeleton';
import { Button } from 'primeng/button';
import { Card } from 'primeng/card';
import { NgClass } from '@angular/common';
import { STORE_FRONT_RESERVATION_ROUTE } from '@store/store.routes';
import { DATES_ROUTE } from '@store/pages/reservation/reservation.routes';
import { Avatar } from 'primeng/avatar';
import { RadioButton } from 'primeng/radiobutton';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-staff',
  imports: [
    Message,
    Skeleton,
    Button,
    Card,
    NgClass,
    Avatar,
    RadioButton,
    FormsModule
  ],
  templateUrl: './staff.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffComponent {
  constructor() {
    const s = this.service1.reservationState().staff;
    if (s) this.staff.set(s);
  }

  private readonly service1 = inject(ReservationService);
  private readonly service2 = inject(StaffService);
  private readonly router = inject(Router);

  protected readonly state = ApiState;
  protected readonly staffs = toSignal(this.service2.staffs(), {
    initialValue: { state: ApiState.LOADING, data: [] } as ApiResponse<
      StaffModel[]
    >
  });

  protected readonly staff = signal<StaffModel | undefined>(undefined);

  protected readonly submit = async () => {
    const staff = this.staff();
    if (!staff) return;
    this.service1.setStaff(staff);
    await this.router.navigate([
      `/${STORE_FRONT_RESERVATION_ROUTE}/${DATES_ROUTE}`
    ]);
  };
}

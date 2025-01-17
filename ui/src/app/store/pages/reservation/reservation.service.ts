import { Injectable, signal } from '@angular/core';
import { ReservationModel } from './reservation.model';
import { ServiceTypeModel } from './service-type/service-type.model';
import { StaffModel } from './staff/staff.model';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  readonly reservationState = signal<ReservationModel>({
    services: undefined,
    staff: undefined
  });

  readonly setServiceTypes = (services: ServiceTypeModel[]) =>
    this.reservationState.set({
      services: services,
      staff: undefined
    });

  readonly setStaff = (staff: StaffModel) =>
    this.reservationState.set({
      services: this.reservationState().services,
      staff: staff
    });
}

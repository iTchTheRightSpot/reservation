import { Injectable, signal } from '@angular/core';
import { ReservationModel } from './reservation.model';
import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { StaffModel } from '@shared/model/shared.model';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  readonly reservationState = signal<ReservationModel>({
    services: undefined,
    staff: undefined,
    datetime: undefined
  });

  readonly setServiceTypes = (services: ServiceTypeModel[]) =>
    this.reservationState.set({
      services: services,
      staff: undefined,
      datetime: undefined
    });

  readonly setStaff = (staff: StaffModel) =>
    this.reservationState.set({
      services: this.reservationState().services,
      staff: staff,
      datetime: undefined
    });

  readonly setDateTime = (date: string | undefined) => {
    const model = this.reservationState();
    this.reservationState.set({
      services: model.services,
      staff: model.staff,
      datetime: date
    });
  };

  readonly clear = () =>
    this.reservationState.set({
      services: undefined,
      staff: undefined,
      datetime: undefined
    });
}

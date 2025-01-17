import { Router } from '@angular/router';
import { inject } from '@angular/core';
import { ReservationService } from './reservation.service';
import { STORE_FRONT_RESERVATION_ROUTE } from '@store/store.routes';
import { SERVICE_TYPE_ROUTE, STAFF_ROUTE } from './reservation.routes';

export const reservationStaffGuard = async () => {
  const obj = inject(ReservationService).reservationState().services;

  if (!obj || obj.length < 1) {
    await inject(Router).navigate([
      `${STORE_FRONT_RESERVATION_ROUTE}/${SERVICE_TYPE_ROUTE}`
    ]);
    return false;
  }

  return true;
};

export const datesGuard = async () => {
  const service = inject(ReservationService).reservationState();
  const bool =
    !service.services || service.services.length < 1 || !service.staff;

  if (bool) {
    await inject(Router).navigate([
      `${STORE_FRONT_RESERVATION_ROUTE}/${STAFF_ROUTE}`
    ]);
    return false;
  }

  return true;
};

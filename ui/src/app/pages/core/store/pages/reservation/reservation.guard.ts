import { Router } from '@angular/router';
import { inject } from '@angular/core';
import { ReservationService } from './reservation.service';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';
import { RESERVATION_ROUTE } from '@store/store.routes';
import {
  CONFIRM_ROUTE,
  SERVICE_TYPE_ROUTE,
  STAFF_ROUTE
} from './reservation.routes';

export const reservationStaffGuard = async () => {
  const obj = inject(ReservationService).reservationState().services;

  if (!obj || obj.length < 1) {
    await inject(Router).navigate([
      `${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${SERVICE_TYPE_ROUTE}`
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
    await inject(Router).navigate([`${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${STAFF_ROUTE}`]);
    return false;
  }

  return true;
};

export const confirmGuard = async () => {
  const service = inject(ReservationService).reservationState();
  const bool =
    !service.services ||
    service.services.length < 1 ||
    !service.staff ||
    !service.dateTime;

  if (bool) {
    await inject(Router).navigate([`${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${CONFIRM_ROUTE}`]);
    return false;
  }

  return true;
};

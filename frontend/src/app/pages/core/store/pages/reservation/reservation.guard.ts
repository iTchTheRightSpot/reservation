import { Router } from '@angular/router';
import { inject } from '@angular/core';
import { ReservationService } from './reservation.service';
import { RootRoutes } from '@root/app.routes';
import { CoreRoutes } from '@pages/core/core.routes';
import { StoreRoutes } from '@store/store.routes';
import { ReservationRoutes } from './reservation.routes';

export const reservationStaffGuard = async () => {
  const obj = inject(ReservationService).reservationState().services;

  if (!obj || obj.length < 1) {
    await inject(Router).navigate([
      `${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.RESERVATION}/${ReservationRoutes.SERVICES}`
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
      `${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.RESERVATION}/${ReservationRoutes.STAFFS}`
    ]);
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
    !service.datetime;

  if (bool) {
    await inject(Router).navigate([
      `${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.RESERVATION}/${ReservationRoutes.DATES}`
    ]);
    return false;
  }

  return true;
};

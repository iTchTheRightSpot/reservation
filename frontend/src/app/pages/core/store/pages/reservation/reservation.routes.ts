import { Routes } from '@angular/router';
import {
  confirmGuard,
  datesGuard,
  reservationStaffGuard
} from './reservation.guard';

export const ReservationRoutes = {
  SERVICES: 'services',
  STAFFS: 'staffs',
  DATES: 'dates',
  CONFIRM: 'confirm'
};

export const routes: Routes = [
  {
    path: ReservationRoutes.SERVICES,
    loadComponent: () =>
      import('./pages/service-type/service-type.component').then(
        m => m.ServiceTypeComponent
      )
  },
  {
    path: ReservationRoutes.STAFFS,
    canActivate: [reservationStaffGuard],
    loadComponent: () =>
      import('./pages/staff/staff.component').then(m => m.StaffComponent)
  },
  {
    path: ReservationRoutes.DATES,
    canActivate: [datesGuard],
    loadComponent: () =>
      import('./pages/dates/dates.component').then(m => m.DatesComponent)
  },
  {
    path: ReservationRoutes.CONFIRM,
    canActivate: [confirmGuard],
    loadComponent: () =>
      import('./pages/confirm/confirm.component').then(m => m.ConfirmComponent)
  },
  {
    path: '',
    redirectTo: ReservationRoutes.SERVICES,
    pathMatch: 'full'
  }
];

import { Routes } from '@angular/router';
import {
  confirmGuard,
  datesGuard,
  reservationStaffGuard
} from './reservation.guard';

export const SERVICE_TYPES_ROUTE = 'services';
export const STAFFS_ROUTE = 'staffs';
export const DATES_ROUTE = 'dates';
export const CONFIRM_ROUTE = 'confirm';

export const routes: Routes = [
  {
    path: SERVICE_TYPES_ROUTE,
    loadComponent: () =>
      import(
        '@store/pages/reservation/pages/service-type/service-type.component'
      ).then(m => m.ServiceTypeComponent)
  },
  {
    path: STAFFS_ROUTE,
    canActivate: [reservationStaffGuard],
    loadComponent: () =>
      import('@store/pages/reservation/pages/staff/staff.component').then(
        m => m.StaffComponent
      )
  },
  {
    path: DATES_ROUTE,
    canActivate: [datesGuard],
    loadComponent: () =>
      import('@store/pages/reservation/pages/dates/dates.component').then(
        m => m.DatesComponent
      )
  },
  {
    path: CONFIRM_ROUTE,
    canActivate: [confirmGuard],
    loadComponent: () =>
      import('@store/pages/reservation/pages/confirm/confirm.component').then(
        m => m.ConfirmComponent
      )
  },
  {
    path: '',
    redirectTo: SERVICE_TYPES_ROUTE,
    pathMatch: 'full'
  }
];

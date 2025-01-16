import { Routes } from '@angular/router';

export const SERVICE_TYPE_ROUTE = 'service';
export const STAFF_ROUTE = 'staff';

export const routes: Routes = [
  {
    path: SERVICE_TYPE_ROUTE,
    loadComponent: () =>
      import(
        '@store/pages/reservation/service-type/service-type.component'
      ).then(m => m.ServiceTypeComponent)
  },
  {
    path: STAFF_ROUTE,
    loadComponent: () =>
      import('@store/pages/reservation/staff/staff.component').then(
        m => m.StaffComponent
      )
  },
  {
    path: '',
    redirectTo: SERVICE_TYPE_ROUTE,
    pathMatch: 'full'
  }
];
